package service

import (
	"fmt"
	"net/http"

	"github.com/haiyiyun/plugins/upload/database/model"
	"github.com/haiyiyun/plugins/upload/service/base"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Upload interface {
	SetUserID(userID primitive.ObjectID)
	SetUserIDFromRequestClaims(r *http.Request) error
	SaveEncodeFile(*http.Request, string, string) (*model.Upload, error)
	SaveFormFile(*http.Request, string, string) (*model.Upload, error)
	DeleteFile(fileID primitive.ObjectID, userID primitive.ObjectID) error // 添加用户ID参数
}

var storageDrivers = make(map[string]func(*base.Service) (Upload, error))

func RegisterStorage(name string, initFunc func(*base.Service) (Upload, error)) {
	storageDrivers[name] = initFunc
}

func NewUploadService(baseService *base.Service) (Upload, error) {
	driverName := baseService.Config.StorageType
	if initFunc, ok := storageDrivers[driverName]; ok {
		return initFunc(baseService)
	}
	return nil, fmt.Errorf("不支持的存储类型: %s", driverName)
}
