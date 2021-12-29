package local

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/haiyiyun/plugins/upload/database/model"
	"github.com/haiyiyun/plugins/upload/database/model/upload"
	"github.com/haiyiyun/plugins/upload/predefined"
	"github.com/haiyiyun/utils/http/request"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (self *Service) SetUserID(userID primitive.ObjectID) {
	self.userID = userID
}

func (self *Service) SetUserIDFromRequestClaims(r *http.Request) error {
	if claims := request.GetClaims(r); claims != nil {
		if uid, err := primitive.ObjectIDFromHex(claims.Issuer); err == nil {
			self.userID = uid
			return nil
		} else {
			return errors.New(predefined.ErrorNotFoundUserIDFromRequestClaims)
		}
	} else {
		return errors.New(predefined.ErrorNotFoundClaimsFromRequest)
	}
}

func (self *Service) generateUploadName(contentType, originalFileName string) (fileType, fileDir, fileName, fileExt string, err error) {
	if contentTypes := strings.Split(contentType, "/"); len(contentTypes) > 1 {
		if originalFileName == "" {
			fileExt = "." + contentTypes[1]
		} else {
			fileExt = filepath.Ext(originalFileName)
		}

		//TODO 后缀名过滤

		//类型过滤
		fileType = predefined.UploadTypeFile
		pathType := self.Config.UploadFileDirectory
		switch contentTypes[0] {
		case "image":
			fileType = predefined.UploadTypeImage
			pathType = self.Config.UploadImageDirectory
		case "video", "audio":
			switch contentTypes[1] {
			case "x-ms-asf", "avi", "x-ivf", "x-mpeg", "mp4", "mpeg4", "x-sgi-movie", "mpeg", "x-mpg", "mpg", "vnd.rn-realvideo", "x-ms-wm", "x-ms-wmv", "x-ms-wmx", "x-mei-aac",
				"aiff", "basic", "x-liquid-file", "x-liquid-secure", "x-la-lms", "mpegurl", "mid", "x-musicnet-download", "x-musicnet-stream", "mp1", "mp2", "mp3", "rn-mpeg", "scpls",
				"vnd.rn-realaudio", "x-pn-realaudio", "x-pn-realaudio-plugin", "wav", "x-ms-wax", "x-ms-wma":
				fileType = predefined.UploadTypeMedia
				pathType = self.Config.UploadMediaDirectory
			}
		case "application":
			switch contentTypes[1] {
			case "msword",
				"vnd.ms-excel",
				"vnd.ms-powerpoint",
				"vnd.ms-visio.drawing",
				"vnd.openxmlformats-officedocument.wordprocessingml.document",
				"vnd.openxmlformats-officedocument.spreadsheetml.sheet",
				"vnd.openxmlformats-officedocument.presentationml.presentation",
				"pdf":
				fileType = predefined.UploadTypeDocument
				pathType = self.Config.UploadDocumentDirectory
			case "x-zip-compressed":
				fileType = predefined.UploadTypeCompression
				pathType = self.Config.UploadCompressionDirectory
			case "octet-stream":
				switch fileExt {
				case ".xmind", ".rp":
					fileType = predefined.UploadTypeDocument
					pathType = self.Config.UploadDocumentDirectory
				case ".rar", ".gz", ".bz2":
					fileType = predefined.UploadTypeCompression
					pathType = self.Config.UploadCompressionDirectory
				}
			}
		}

		now := time.Now()
		fileDir = filepath.Clean(fmt.Sprintf("%s%s/%d/%d/%d", self.Config.UploadDirectory, pathType, now.Year(), now.Month(), now.Day()))
		if _, err = os.Stat(fileDir); err != nil {
			err = os.MkdirAll(fileDir, 0755)
		}

		random := rand.New(rand.NewSource(now.UnixNano()))
		fileName = fmt.Sprintf("%v", random.Uint64())
	} else {
		err = errors.New("parse Content-Type failed")
	}

	return
}

func (self *Service) relativePath(sPath string) (newPath string) {
	uploadDir := filepath.Clean(self.Config.UploadDirectory)
	newPath = sPath
	if iPos := strings.Index(sPath, uploadDir); iPos > -1 {
		newPath = sPath[iPos+len(uploadDir):]
	}

	return strings.TrimLeft(newPath, "/")
}

func (self *Service) saveFileToPath(contentType, originalFileName string, fileData []byte, appendFileExt bool) (*model.Upload, error) {
	err := errors.New(predefined.ErrorFalidSaveFile)
	if fileType, fileDir, fileName, fileExt, fileErr := self.generateUploadName(contentType, originalFileName); fileErr == nil {
		if originalFileName == "" {
			originalFileName = fileName + fileExt
		}

		if appendFileExt {
			fileName = fileName + fileExt
		}

		filePath := fileDir + "/" + fileName
		if err = ioutil.WriteFile(filePath, fileData, 0755); err == nil {
			var fileinfo os.FileInfo
			if fileinfo, err = os.Stat(filePath); err == nil {
				mf := &model.Upload{
					ID:               primitive.NewObjectID(),
					Type:             fileType,
					Storage:          predefined.UploadStorageLocal,
					UserID:           self.userID,
					ContentType:      contentType,
					OriginalFileName: originalFileName,
					FileName:         fileName,
					FileExt:          fileExt,
					Path:             self.relativePath(filePath),
					URL:              self.relativePath(filePath),
					Size:             fileinfo.Size(),
				}

				return mf, nil
			}
		}
	} else {
		err = fileErr
	}

	return nil, err
}

func (self *Service) encodeDataToFile(encode_string string, appendFileExt bool) (*model.Upload, error) {
	err := errors.New(predefined.ErrorNotFoundEncodeData)
	if encode_strings := strings.Split(encode_string, ";"); len(encode_strings) > 1 {
		if encode_datas := strings.Split(encode_strings[1], ","); len(encode_datas) > 1 {
			encode_type := encode_datas[0]
			encode_data := encode_datas[1]
			switch encode_type {
			case "base64":
				var file_data []byte
				if file_data, err = base64.StdEncoding.DecodeString(encode_data); err == nil {
					return self.saveFileToPath(encode_strings[0], "", file_data, appendFileExt)
				}
			}
		}
	}

	return nil, err
}

func (self *Service) saveFormFile(r *http.Request, fileFormName string, bEncode bool, remark string) (fm *model.Upload, err error) {
	if bEncode {
		form_value := r.FormValue(fileFormName)
		if form_value != "" {
			if datas := strings.Split(form_value, ":"); len(datas) > 1 {
				fm, err = self.encodeDataToFile(datas[1], self.Config.AppendFileExt)
			}
		} else {
			err = errors.New(predefined.ErrorNotFoundFormData)
		}
	} else {
		if f, fh, fErr := r.FormFile(fileFormName); fErr != nil {
			err = fErr
		} else {
			if fileData, ioErr := ioutil.ReadAll(f); ioErr != nil {
				err = ioErr
			} else {
				fm, err = self.saveFileToPath(fh.Header.Get("Content-Type"), fh.Filename, fileData, self.Config.AppendFileExt)
			}
		}
	}

	if err == nil && fm != nil {
		uploadModel := upload.NewModel(self.M)
		if remark != "" {
			fm.Remark = remark
		}

		_, err = uploadModel.Create(context.Background(), fm)
	}

	return
}

func (self *Service) SaveEncodeFile(r *http.Request, fileFormName, remark string) (fm *model.Upload, err error) {
	if !self.AllowUploadLocal {
		err = errors.New(predefined.ErrorNotAllowUploadLocal)
		return
	}

	return self.saveFormFile(r, fileFormName, true, remark)
}

func (self *Service) SaveFormFile(r *http.Request, fileFormName, remark string) (fm *model.Upload, err error) {
	if !self.AllowUploadLocal {
		err = errors.New(predefined.ErrorNotAllowUploadLocal)
		return
	}

	return self.saveFormFile(r, fileFormName, false, remark)
}
