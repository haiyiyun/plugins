package aliyun

import (
	"context"
	"errors"
	"net/http"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/haiyiyun/plugins/upload/database/model"
	"github.com/haiyiyun/plugins/upload/database/model/upload"
	"github.com/haiyiyun/plugins/upload/predefined"
	"github.com/haiyiyun/plugins/upload/service"
	"github.com/haiyiyun/plugins/upload/service/base"
	"github.com/haiyiyun/utils/http/request"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	service.RegisterStorage("aliyun", NewService)
}

type Service struct {
	*base.Service
	client *oss.Client
	bucket *oss.Bucket
	userID primitive.ObjectID
}

func NewService(baseService *base.Service) (service.Upload, error) {
	client, err := oss.New(baseService.Config.AliyunConfig.AliyunEndpoint, baseService.Config.AliyunConfig.AliyunAccessKeyID, baseService.Config.AliyunConfig.AliyunAccessKeySecret)
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(baseService.Config.AliyunConfig.AliyunBucketName)
	if err != nil {
		return nil, err
	}

	return &Service{
		Service: baseService,
		client:  client,
		bucket:  bucket, // 添加这行，将bucket存储到结构体
	}, nil
}

func (s *Service) SetUserID(userID primitive.ObjectID) {
	s.userID = userID
}

func (s *Service) SetUserIDFromRequestClaims(r *http.Request) error {
	if claims := request.GetClaims(r); claims != nil {
		if uid, err := primitive.ObjectIDFromHex(claims.Issuer); err == nil {
			s.userID = uid
			return nil
		} else {
			return errors.New(predefined.ErrorNotFoundUserIDFromRequestClaims)
		}
	} else {
		return errors.New(predefined.ErrorNotFoundClaimsFromRequest)
	}
}

func (s *Service) SaveEncodeFile(r *http.Request, fileFormName, remark string) (fm *model.Upload, err error) {
	if s.Config.AliyunConfig.AliyunDisableUpload {
		return nil, errors.New(predefined.ErrorNotAllowUploadAliyun)
	}
	return s.saveFormFile(r, fileFormName, true, remark)
}

func (s *Service) SaveFormFile(r *http.Request, fileFormName, remark string) (fm *model.Upload, err error) {
	if s.Config.AliyunConfig.AliyunDisableUpload {
		return nil, errors.New(predefined.ErrorNotAllowUploadAliyun)
	}
	return s.saveFormFile(r, fileFormName, false, remark)
}

func (s *Service) DeleteFile(fileID primitive.ObjectID, userID primitive.ObjectID) error {
	if s.Config.AliyunDisableBucketCRUD {
		return errors.New("aliyun bucket CRUD disabled")
	}

	uploadModel := upload.NewModel(s.M)

	// 构建查询条件：文件ID + 用户ID
	filter := uploadModel.FilterByID(fileID)
	if !userID.IsZero() {
		filter = append(filter, uploadModel.FilterByUserID(userID)...)
	}

	// 查找文件记录
	sr := uploadModel.FindOne(context.Background(), filter)
	var uploadFile model.Upload
	if err := sr.Decode(&uploadFile); err != nil {
		return err
	}

	// 删除云存储文件
	if err := s.bucket.DeleteObject(uploadFile.Path); err != nil {
		return err
	}

	// 删除数据库记录
	_, err := uploadModel.DeleteOne(context.Background(), filter)
	return err
}
