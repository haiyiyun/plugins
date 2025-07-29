package upload

import (
	"net/http"
	"os"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/plugins/upload/database/model"
	"github.com/haiyiyun/plugins/upload/database/model/upload"
	"github.com/haiyiyun/plugins/upload/predefined"
	"github.com/haiyiyun/utils/help"
	"github.com/haiyiyun/utils/http/request"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/utils/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (self *Service) Route_GET_File(rw http.ResponseWriter, r *http.Request) {
	var requestUID predefined.RequestServeUploadID
	if err := validator.FormStruct(&requestUID, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, "无效的文件ID参数: "+err.Error())
		return
	}

	uploadModel := upload.NewModel(self.M)
	sr := uploadModel.FindOne(r.Context(), uploadModel.FilterByID(requestUID.ID))
	var uploadFile model.Upload
	if err := sr.Decode(&uploadFile); err != nil {
		response.JSON(rw, http.StatusNotFound, nil, "文件不存在或已被删除")
		return
	} else {
		switch uploadFile.Storage {
		case predefined.UploadStorageLocal:
			if self.AllowDownloadLocal {
				fullPath := self.Config.UploadDirectory + uploadFile.Path
				if _, err := os.Stat(fullPath); os.IsNotExist(err) {
					response.JSON(rw, http.StatusNotFound, nil, "本地文件不存在")
					return
				}
				http.ServeFile(rw, r, fullPath)
			} else {
				response.JSON(rw, http.StatusServiceUnavailable, nil, "本地文件下载功能已禁用")
			}
		case predefined.UploadStorageAliyun:
			// 阿里云OSS下载处理
			if self.Config.AliyunDisableDownload {
				response.JSON(rw, http.StatusServiceUnavailable, nil, "阿里云下载功能已禁用")
				return
			}
			http.Redirect(rw, r, uploadFile.URL, http.StatusFound)
		case predefined.UploadStorageTencent:
			// 腾讯云COS下载处理
			if self.Config.TencentDisableDownload {
				response.JSON(rw, http.StatusServiceUnavailable, nil, "腾讯云下载功能已禁用")
				return
			}
			http.Redirect(rw, r, uploadFile.URL, http.StatusFound)
		case predefined.UploadStorageQiniu:
			// 七牛云下载处理
			if self.Config.QiniuDisableDownload { // 使用正确的下载禁用标志
				response.JSON(rw, http.StatusServiceUnavailable, nil, "七牛云下载功能已禁用")
				return
			}
			http.Redirect(rw, r, uploadFile.URL, http.StatusFound)
		default:
			response.JSON(rw, http.StatusServiceUnavailable, nil, "不支持的存储类型: "+uploadFile.Storage)
		}
	}
}

func (self *Service) Route_POST_File(rw http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(self.Config.MaxUploadFileSize)

	var requestF predefined.RequestServeFile
	if err := validator.FormStruct(&requestF, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	// 设置用户ID
	if self.Config.CheckUser {
		if err := self.uploadService.SetUserIDFromRequestClaims(r); err != nil {
			response.JSON(rw, http.StatusUnauthorized, nil, "")
			return
		}
	}

	fileBase64Data := r.FormValue(predefined.FormNameFileBase64Data)
	var fm *model.Upload
	var err error

	if fileBase64Data != "" {
		if requestF.FileType == predefined.UploadTypeImage {
			fm, err = self.uploadService.SaveEncodeFile(r, predefined.FormNameFileBase64Data, requestF.Remark)
		} else {
			response.JSON(rw, http.StatusBadRequest, nil, "不支持的文件类型")
			return
		}
	} else {
		switch requestF.FileType {
		case predefined.UploadTypeImage,
			predefined.UploadTypeMedia,
			predefined.UploadTypeDocument,
			predefined.UploadTypeCompression,
			predefined.UploadTypeFile:
			fm, err = self.uploadService.SaveFormFile(r, predefined.FormNameFile, requestF.Remark)
		default:
			response.JSON(rw, http.StatusBadRequest, nil, "不支持的文件类型")
			return
		}
	}

	if err != nil {
		log.Error(err)
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
		return
	}

	result := help.M{
		"id": fm.ID.Hex(),
	}

	// 根据存储类型生成下载URL
	if self.Config.PublishDownloadUrl {
		switch fm.Storage {
		case predefined.UploadStorageLocal:
			result["url"] = self.Config.DownloadLocalUrlDirectory + fm.URL
		case predefined.UploadStorageAliyun,
			predefined.UploadStorageTencent,
			predefined.UploadStorageQiniu:
			result["url"] = fm.URL // 云存储直接使用返回的URL
		}
	}

	response.JSON(rw, 0, result, "上传成功")
}

func (self *Service) Route_DELETE_File(rw http.ResponseWriter, r *http.Request) {
	userID := primitive.NilObjectID

	if self.Config.CheckUser {
		if cliaims := request.GetClaims(r); cliaims != nil {
			if uid, err := primitive.ObjectIDFromHex(cliaims.Issuer); err != nil {
				response.JSON(rw, http.StatusUnauthorized, nil, "")
				return
			} else {
				userID = uid
			}
		} else {
			response.JSON(rw, http.StatusUnauthorized, nil, "")
			return
		}
	}

	var requestUID predefined.RequestServeUploadID
	if err := validator.FormStruct(&requestUID, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	// 调用删除方法，传递文件ID和用户ID
	if err := self.uploadService.DeleteFile(requestUID.ID, userID); err != nil {
		log.Error("Delete file failed:", err)
		response.JSON(rw, http.StatusInternalServerError, nil, "文件删除失败")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}
