package upload

import (
	"net/http"
	"os"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/plugins/upload/database/model"
	"github.com/haiyiyun/plugins/upload/database/model/upload"
	"github.com/haiyiyun/plugins/upload/predefined"
	"github.com/haiyiyun/plugins/upload/service/local"
	"github.com/haiyiyun/utils/help"
	"github.com/haiyiyun/utils/http/request"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (self *Service) Route_GET_File(rw http.ResponseWriter, r *http.Request) {
	fileIDStr := r.URL.Query().Get("id")

	valid := &validator.Validation{}
	valid.Required(fileIDStr).Key("id").Message("id字段为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	fileID, fileIDErr := primitive.ObjectIDFromHex(fileIDStr)
	if fileIDErr != nil {
		response.JSON(rw, http.StatusBadRequest, nil, "")
		return
	}

	uploadModel := upload.NewModel(self.M)
	sr := uploadModel.FindOne(r.Context(), uploadModel.FilterByID(fileID))
	var uploadFile model.Upload
	if err := sr.Decode(&uploadFile); err != nil {
		response.JSON(rw, http.StatusNotFound, nil, "")
		return
	} else {
		if uploadFile.Storage == predefined.UploadStorageLocal {
			if self.AllowDownloadLocal {
				http.ServeFile(rw, r, self.Config.UploadDirectory+uploadFile.Path)
			} else {
				response.JSON(rw, http.StatusServiceUnavailable, nil, "")
			}
		} else {
			response.JSON(rw, http.StatusServiceUnavailable, nil, "")
		}
	}
}

func (self *Service) Route_POST_File(rw http.ResponseWriter, r *http.Request) {
	fileType := r.FormValue("file_type")
	fileBase64Data := r.FormValue(predefined.FormNameFileBase64Data)
	remark := r.FormValue("remark")

	valid := &validator.Validation{}
	valid.Required(fileType).Key("file_type").Message("file_type字段为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	if fileBase64Data != "" {
		if fileType == predefined.UploadTypeImage {
			uploadLocal := local.NewService(self.Service.Service)
			setUIDErr := uploadLocal.SetUserIDFromRequestClaims(r)
			if self.Config.CheckUser {
				if setUIDErr != nil {
					response.JSON(rw, http.StatusUnauthorized, nil, "")
					return
				}
			}

			if fm, err := uploadLocal.SaveEncodeFile(r, predefined.FormNameFileBase64Data, remark); err != nil {
				log.Error(err)

				if err.Error() == predefined.ErrorNotAllowUploadLocal {
					response.JSON(rw, http.StatusServiceUnavailable, nil, "")
					return
				}
			} else {
				response.JSON(rw, 0, help.M{
					"id":  fm.ID.Hex(),
					"url": self.Config.DownloadLocalUrlDirectory + fm.URL,
				}, "上传成功")

				return
			}
		}
	} else {
		switch fileType {
		case predefined.UploadTypeImage, predefined.UploadTypeMedia, predefined.UploadTypeFile, "":
			uploadLocal := local.NewService(self.Service.Service)
			setUIDErr := uploadLocal.SetUserIDFromRequestClaims(r)
			if self.Config.CheckUser {
				if setUIDErr != nil {
					response.JSON(rw, http.StatusUnauthorized, nil, "")
					return
				}
			}

			if fm, err := uploadLocal.SaveFormFile(r, predefined.FormNameFile, remark); err != nil {
				log.Error(err)

				if err.Error() == predefined.ErrorNotAllowUploadLocal {
					response.JSON(rw, http.StatusServiceUnavailable, nil, "")
					return
				}
			} else {
				response.JSON(rw, 0, help.M{
					"id":  fm.ID.Hex(),
					"url": self.Config.DownloadLocalUrlDirectory + fm.URL,
				}, "上传成功")

				return
			}
		}
	}

	response.JSON(rw, http.StatusBadRequest, nil, "上传失败")
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

	fileIDStr := r.URL.Query().Get("id")

	valid := &validator.Validation{}
	valid.Required(fileIDStr).Key("id").Message("id字段为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	fileID, fileIDErr := primitive.ObjectIDFromHex(fileIDStr)
	if fileIDErr != nil {
		response.JSON(rw, http.StatusBadRequest, nil, "")
		return
	}

	uploadModel := upload.NewModel(self.M)
	filter := uploadModel.FilterByID(fileID)
	filter = append(filter, uploadModel.FilterByUserID(userID)...)

	sr := uploadModel.FindOne(r.Context(), uploadModel.FilterByID(fileID))
	var uploadFile model.Upload
	if err := sr.Decode(&uploadFile); err != nil {
		response.JSON(rw, http.StatusNotFound, nil, "")
	} else {
		if uploadFile.Storage == predefined.UploadStorageLocal {
			filePath := self.Config.UploadDirectory + uploadFile.Path
			if err := os.Remove(filePath); err != nil {
				response.JSON(rw, http.StatusNotFound, nil, "")
			} else {
				if dr, err := uploadModel.DeleteOne(r.Context(), filter); err != nil || dr.DeletedCount == 0 {
					response.JSON(rw, http.StatusNotFound, nil, "")
				} else {
					response.JSON(rw, 0, nil, "")
				}
			}
		} else {
			response.JSON(rw, http.StatusServiceUnavailable, nil, "")
		}
	}
}
