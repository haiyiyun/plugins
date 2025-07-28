package qiniuyun

import (
	"context"
	"errors"
	"net/http"

	"github.com/haiyiyun/plugins/upload/database/model"
	"github.com/haiyiyun/plugins/upload/database/model/upload"
	"github.com/haiyiyun/plugins/upload/predefined"
	"github.com/haiyiyun/plugins/upload/service"
	"github.com/haiyiyun/plugins/upload/service/base"
	"github.com/haiyiyun/utils/http/request"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	*base.Service
	userID   primitive.ObjectID
	upToken  string
	upConfig storage.Config
	mac      *auth.Credentials // 添加 mac 字段
}

func NewService(s *base.Service) (service.Upload, error) {
	qiniuConfig := s.Config.QiniuConfig

	// 构建七牛云配置
	cfg := storage.Config{}
	// 是否使用HTTPS
	cfg.UseHTTPS = qiniuConfig.QiniuUseHTTPS
	// 是否使用CDN加速上传
	cfg.UseCdnDomains = qiniuConfig.QiniuUseCdnDomains

	// 初始化凭证
	mac := auth.New(qiniuConfig.QiniuAccessKey, qiniuConfig.QiniuSecretKey) // 创建认证凭证
	putPolicy := storage.PutPolicy{
		Scope: qiniuConfig.QiniuBucketName,
	}
	upToken := putPolicy.UploadToken(mac)

	return &Service{
		Service:  s,
		upToken:  upToken,
		upConfig: cfg,
		mac:      mac, // 保存认证凭证
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
	if s.Config.QiniuConfig.QiniuDisableUpload {
		return nil, errors.New(predefined.ErrorNotAllowUploadQiniu)
	}
	return s.saveFormFile(r, fileFormName, true, remark)
}

func (s *Service) SaveFormFile(r *http.Request, fileFormName, remark string) (fm *model.Upload, err error) {
	if s.Config.QiniuConfig.QiniuDisableUpload {
		return nil, errors.New(predefined.ErrorNotAllowUploadQiniu)
	}
	return s.saveFormFile(r, fileFormName, false, remark)
}

func (s *Service) DeleteFile(fileID primitive.ObjectID, userID primitive.ObjectID) error {
	if s.Config.QiniuConfig.QiniuDisableUpload {
		return errors.New("qiniu file operations disabled")
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
	bucketManager := storage.NewBucketManager(s.mac, &s.upConfig)
	if err := bucketManager.Delete(s.Config.QiniuBucketName, uploadFile.Path); err != nil {
		return err
	}

	// 删除数据库记录
	_, err := uploadModel.DeleteOne(context.Background(), filter)
	return err
}
