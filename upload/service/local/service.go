package local

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/haiyiyun/plugins/upload/database/model"
	"github.com/haiyiyun/plugins/upload/database/model/upload"
	"github.com/haiyiyun/plugins/upload/predefined"
	"github.com/haiyiyun/plugins/upload/service"
	"github.com/haiyiyun/plugins/upload/service/base"
	"github.com/haiyiyun/utils/http/request"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	*base.Service
	userID primitive.ObjectID
}

func NewService(s *base.Service) (service.Upload, error) {
	return &Service{
		Service: s,
		userID:  primitive.NilObjectID,
	}, nil
}

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

func (self *Service) DeleteFile(fileID primitive.ObjectID, userID primitive.ObjectID) error {
	uploadModel := upload.NewModel(self.M)

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

	// 删除物理文件
	fullPath := filepath.Join(self.Config.UploadDirectory, uploadFile.Path)
	if err := os.Remove(fullPath); err != nil {
		return err
	}

	// 删除数据库记录
	_, err := uploadModel.DeleteOne(context.Background(), filter)
	return err
}

func init() {
	service.RegisterStorage("local", NewService)
}
