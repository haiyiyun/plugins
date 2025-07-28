package tencent

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/haiyiyun/plugins/upload/database/model"
	"github.com/haiyiyun/plugins/upload/database/model/upload"
	"github.com/haiyiyun/plugins/upload/predefined"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *Service) generateObjectKey(contentType, originalFileName string) (fileType, objectKey, fileExt string, err error) {
	if contentTypes := strings.Split(contentType, "/"); len(contentTypes) > 1 {
		if originalFileName == "" {
			fileExt = "." + contentTypes[1]
		} else {
			fileExt = filepath.Ext(originalFileName)
		}

		// 类型过滤
		fileType = predefined.UploadTypeFile
		switch contentTypes[0] {
		case "image":
			fileType = predefined.UploadTypeImage
		case "video", "audio":
			fileType = predefined.UploadTypeMedia
		case "application":
			switch contentTypes[1] {
			case "msword", "vnd.ms-excel", "vnd.ms-powerpoint", "vnd.ms-visio.drawing",
				"vnd.openxmlformats-officedocument.wordprocessingml.document",
				"vnd.openxmlformats-officedocument.spreadsheetml.sheet",
				"vnd.openxmlformats-officedocument.presentationml.presentation", "pdf":
				fileType = predefined.UploadTypeDocument
			case "x-zip-compressed":
				fileType = predefined.UploadTypeCompression
			case "octet-stream":
				switch fileExt {
				case ".xmind", ".rp":
					fileType = predefined.UploadTypeDocument
				case ".rar", ".gz", ".bz2":
					fileType = predefined.UploadTypeCompression
				}
			}
		}

		now := time.Now()
		random := rand.New(rand.NewSource(now.UnixNano()))
		fileName := fmt.Sprintf("%v", random.Uint64())
		if s.Config.BaseConfig.AppendFileExt {
			fileName += fileExt
		}

		// 生成对象键：类型/年/月/日/文件名
		objectKey = fmt.Sprintf("%s/%d/%02d/%02d/%s", fileType, now.Year(), now.Month(), now.Day(), fileName)
	} else {
		err = errors.New("parse Content-Type failed")
	}

	return
}

func (s *Service) saveFormFile(r *http.Request, fileFormName string, bEncode bool, remark string) (fm *model.Upload, err error) {
	if bEncode {
		formValue := r.FormValue(fileFormName)
		if formValue != "" {
			if datas := strings.Split(formValue, ":"); len(datas) > 1 {
				// 移除未使用的 appendFileExt 参数
				fm, err = s.encodeDataToFile(datas[1], remark)
			}
		} else {
			err = errors.New(predefined.ErrorNotFoundFormData)
		}
	} else {
		if f, fh, fErr := r.FormFile(fileFormName); fErr != nil {
			err = fErr
		} else {
			defer f.Close()

			contentType := fh.Header.Get("Content-Type")
			originalFileName := fh.Filename
			fileType, objectKey, fileExt, genErr := s.generateObjectKey(contentType, originalFileName)
			if genErr != nil {
				return nil, genErr
			}

			// 使用腾讯云COS SDK上传文件
			_, err = s.client.Object.Put(context.Background(), objectKey, f, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to upload to Tencent COS: %w", err)
			}

			// 获取文件信息
			resp, err := s.client.Object.Head(context.Background(), objectKey, nil)
			if err != nil {
				return nil, fmt.Errorf("failed to get file info from Tencent COS: %w", err)
			}
			fileSize := resp.ContentLength

			// 构建文件URL
			fileURL := s.Config.TencentConfig.TencentBaseURL + objectKey
			if s.Config.TencentConfig.TencentSecure {
				fileURL = "https://" + fileURL
			} else {
				fileURL = "http://" + fileURL
			}

			// 构建文件信息
			fm = &model.Upload{
				ID:               primitive.NewObjectID(),
				Type:             fileType,
				Storage:          predefined.UploadStorageTencent,
				UserID:           s.userID,
				ContentType:      contentType,
				OriginalFileName: originalFileName,
				FileName:         filepath.Base(objectKey),
				FileExt:          fileExt,
				Path:             objectKey,
				URL:              fileURL,
				Size:             fileSize,
				Remark:           remark,
			}
		}
	}

	if err == nil && fm != nil {
		uploadModel := upload.NewModel(s.M)
		if remark != "" {
			fm.Remark = remark
		}

		// 保存到数据库
		_, err = uploadModel.Create(context.Background(), fm)

		// 如果数据库保存失败，尝试删除COS上的文件
		if err != nil {
			_, delErr := s.client.Object.Delete(context.Background(), fm.Path)
			if delErr != nil {
				return nil, fmt.Errorf("database error: %w, and failed to delete COS file: %v", err, delErr)
			}
			return nil, fmt.Errorf("database error: %w, COS file deleted", err)
		}
	}

	return
}

// 移除未使用的 appendFileExt 参数
func (s *Service) encodeDataToFile(encodeString string, remark string) (*model.Upload, error) {
	err := errors.New(predefined.ErrorNotFoundEncodeData)
	if encodeStrings := strings.Split(encodeString, ";"); len(encodeStrings) > 1 {
		if encodeDatas := strings.Split(encodeStrings[1], ","); len(encodeDatas) > 1 {
			encodeType := encodeDatas[0]
			encodeData := encodeDatas[1]
			switch encodeType {
			case "base64":
				var fileData []byte
				if fileData, err = base64.StdEncoding.DecodeString(encodeData); err == nil {
					contentType := encodeStrings[0]
					originalFileName := ""

					// 生成对象键
					fileType, objectKey, fileExt, genErr := s.generateObjectKey(contentType, originalFileName)
					if genErr != nil {
						return nil, genErr
					}

					// 上传到COS
					_, err = s.client.Object.Put(context.Background(), objectKey, strings.NewReader(string(fileData)), nil)
					if err != nil {
						return nil, fmt.Errorf("failed to upload base64 data to Tencent COS: %w", err)
					}

					// 获取文件信息
					resp, err := s.client.Object.Head(context.Background(), objectKey, nil)
					if err != nil {
						return nil, fmt.Errorf("failed to get base64 file info from Tencent COS: %w", err)
					}
					fileSize := resp.ContentLength

					// 构建文件URL
					fileURL := s.Config.TencentConfig.TencentBaseURL + objectKey
					if s.Config.TencentConfig.TencentSecure {
						fileURL = "https://" + fileURL
					} else {
						fileURL = "http://" + fileURL
					}

					// 返回文件模型
					return &model.Upload{
						ID:               primitive.NewObjectID(),
						Type:             fileType,
						Storage:          predefined.UploadStorageTencent,
						UserID:           s.userID,
						ContentType:      contentType,
						OriginalFileName: originalFileName,
						FileName:         filepath.Base(objectKey),
						FileExt:          fileExt,
						Path:             objectKey,
						URL:              fileURL,
						Size:             fileSize,
						Remark:           remark,
					}, nil
				}
			}
		}
	}
	return nil, err
}
