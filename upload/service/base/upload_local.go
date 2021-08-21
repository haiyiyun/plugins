package base

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

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/utils/http/request"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UploadLocal struct {
	*Service
}

func NewUploadLocal(s *Service) *UploadLocal {
	return &UploadLocal{
		Service: s,
	}
}

func (ufl *UploadLocal) generateUploadName(contentType, originalFileName string) (fileType, fileDir, fileName, fileExt string) {
	if contentTypes := strings.Split(contentType, "/"); len(contentTypes) > 1 {
		if originalFileName == "" {
			fileExt = "." + contentTypes[1]
		} else {
			fileExt = filepath.Ext(originalFileName)
		}

		//TODU 后缀名过滤

		//类型过滤
		fileType = predefined.UploadTypeFile
		pathType := ufl.Config.UploadDirectory
		switch contentTypes[0] {
		case predefined.UploadTypeImage:
			fileType = predefined.UploadTypeImage
			pathType = ufl.Config.UploadImageDirectory
		case predefined.UploadTypeMedia:
			switch fileType {
			case ".mp3", ".wav", ".wma", ".wmv", ".mid", ".avi", ".mpg", ".asf", ".rm", ".rmvb":
				fileType = predefined.UploadTypeMedia
				pathType = ufl.Config.UploadMediaDirectory
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
				pathType = ufl.Config.UploadDocumentDirectory
			case "x-zip-compressed":
				fileType = predefined.UploadTypeCompression
				pathType = ufl.Config.UploadCompressionDirectory
			case "octet-stream":
				switch fileExt {
				case ".xmind", ".rp":
					fileType = predefined.UploadTypeDocument
					pathType = ufl.Config.UploadDocumentDirectory
				case ".rar", ".gz", ".bz2":
					fileType = predefined.UploadTypeCompression
					pathType = ufl.Config.UploadCompressionDirectory
				}
			}
		}

		now := time.Now()
		fileDir = filepath.Clean(fmt.Sprintf("%s/%d/%d/%d", pathType, now.Year(), now.Month(), now.Day()))
		if _, err := os.Stat(fileDir); err != nil {
			if err := os.MkdirAll(fileDir, 0755); err != nil {
				log.Error(err)
			}
		}

		random := rand.New(rand.NewSource(now.UnixNano()))
		fileName = fmt.Sprintf("%v", random.Uint64())
	}

	return
}

func (ufl *UploadLocal) filePathToURL(sPath string) (sURL string) {
	uploadDir := filepath.Clean(ufl.Config.UploadDirectory)
	sURL = sPath
	if iPos := strings.Index(sPath, uploadDir); iPos > -1 {
		sURL = sPath[iPos+len(uploadDir):]
	}

	return
}

func (ufl *UploadLocal) saveFileToPath(contentType, originalFileName string, fileData []byte, appendFileExt bool) (*model.Upload, error) {
	err := errors.New("Fail to save file")
	if fileType, fileDir, fileName, fileExt := ufl.generateUploadName(contentType, originalFileName); fileDir != "" {
		if originalFileName == "" {
			originalFileName = fileName + fileExt
		}

		if appendFileExt {
			fileName = fileName + fileExt
		}

		filePath := fileDir + "/" + fileName
		if err = ioutil.WriteFile(filePath, fileData, 0755); err == nil {
			sURL := ufl.filePathToURL(filePath)
			var fileinfo os.FileInfo
			if fileinfo, err = os.Stat(filePath); err == nil {
				mf := &model.Upload{
					Type:             fileType,
					Storage:          predefined.UploadStorageLocal,
					ContentType:      contentType,
					OriginalFileName: originalFileName,
					FileName:         fileName,
					FileExt:          fileExt,
					Path:             filePath,
					URL:              sURL,
					Size:             fileinfo.Size(),
				}

				return mf, nil
			}
		}
	}

	return nil, err
}

func (ufl *UploadLocal) encodeDataToFile(encode_string string, appendFileExt bool) (*model.Upload, error) {
	err := errors.New("Not Found encode data")
	if encode_strings := strings.Split(encode_string, ";"); len(encode_strings) > 1 {
		if encode_datas := strings.Split(encode_strings[1], ","); len(encode_datas) > 1 {
			encode_type := encode_datas[0]
			encode_data := encode_datas[1]
			switch encode_type {
			case "base64":
				var file_data []byte
				if file_data, err = base64.StdEncoding.DecodeString(encode_data); err == nil {
					return ufl.saveFileToPath(encode_strings[0], "", file_data, appendFileExt)
				}
			}
		}
	}

	return nil, err
}

func (ufl *UploadLocal) saveFormFile(r *http.Request, fileFormName string, bEncode bool) (fm *model.Upload, err error) {
	//TODO 后期改为配置控制
	appendFileExt := true
	if bEncode {
		form_value := r.FormValue(fileFormName)
		if form_value != "" {
			if datas := strings.Split(form_value, ":"); len(datas) > 1 {
				fm, err = ufl.encodeDataToFile(datas[1], appendFileExt)
			}
		} else {
			err = errors.New("Not found form data")
		}
	} else {
		if f, fh, fErr := r.FormFile(fileFormName); fErr != nil {
			err = fErr
		} else {
			if fileData, ioErr := ioutil.ReadAll(f); ioErr != nil {
				err = ioErr
			} else {
				fm, err = ufl.saveFileToPath(fh.Header.Get("Content-Type"), fh.Filename, fileData, appendFileExt)
			}
		}
	}

	if err == nil && fm != nil {
		claims := request.GetClaims(r)
		fm.UserID, _ = primitive.ObjectIDFromHex(claims.Issuer)
		uploadModel := upload.NewModel(ufl.M)
		_, err = uploadModel.Create(context.Background(), fm)
	}

	return
}

func (ufl *UploadLocal) SaveEncodeFile(r *http.Request, fileFormName string) (fm *model.Upload, err error) {
	return ufl.saveFormFile(r, fileFormName, true)
}

func (ufl *UploadLocal) SaveFormFile(r *http.Request, fileFormName string) (fm *model.Upload, err error) {
	return ufl.saveFormFile(r, fileFormName, false)
}
