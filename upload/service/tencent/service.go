package tencent

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/haiyiyun/plugins/upload/database/model"
	"github.com/haiyiyun/plugins/upload/database/model/upload"
	"github.com/haiyiyun/plugins/upload/predefined"
	"github.com/haiyiyun/plugins/upload/service"
	"github.com/haiyiyun/plugins/upload/service/base"
	"github.com/haiyiyun/utils/http/request"
	"github.com/tencentyun/cos-go-sdk-v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	*base.Service
	userID primitive.ObjectID
	client *cos.Client // 腾讯云COS客户端
}

func NewService(s *base.Service) (service.Upload, error) {
	// 从配置中获取腾讯云COS相关参数
	tencentConfig := s.Config.TencentConfig

	// 构建COS的URL
	u, err := url.Parse("https://" + tencentConfig.TencentBucketName + ".cos." + tencentConfig.TencentEndpoint + ".myqcloud.com")
	if err != nil {
		return nil, err
	}

	// 初始化COS客户端
	baseURL := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(baseURL, &http.Client{
		Timeout: 30 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:  tencentConfig.TencentAccessKeyID,
			SecretKey: tencentConfig.TencentAccessKeySecret,
		},
	})

	return &Service{
		Service: s,
		userID:  primitive.NilObjectID,
		client:  client,
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
	if s.Config.TencentConfig.TencentDisableUpload {
		return nil, errors.New(predefined.ErrorNotAllowUploadTencent)
	}
	return s.saveFormFile(r, fileFormName, true, remark)
}

func (s *Service) SaveFormFile(r *http.Request, fileFormName, remark string) (fm *model.Upload, err error) {
	if s.Config.TencentConfig.TencentDisableUpload {
		return nil, errors.New(predefined.ErrorNotAllowUploadTencent)
	}
	return s.saveFormFile(r, fileFormName, false, remark)
}

func (s *Service) DeleteFile(fileID primitive.ObjectID, userID primitive.ObjectID) error {
	if s.Config.TencentDisableBucketCRUD {
		return errors.New("tencent bucket CRUD disabled")
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
	_, err := s.client.Object.Delete(context.Background(), uploadFile.Path)
	if err != nil {
		return err
	}

	// 删除数据库记录
	_, err = uploadModel.DeleteOne(context.Background(), filter)
	return err
}

func init() {
	service.RegisterStorage("tencent", NewService)
}
