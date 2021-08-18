package upload

import (
	"context"
	"flag"

	"github.com/haiyiyun/plugins/upload/database/schema"
	"github.com/haiyiyun/plugins/upload/service"
	serviceUpload "github.com/haiyiyun/plugins/upload/service/upload"

	"github.com/haiyiyun/cache"
	"github.com/haiyiyun/config"
	"github.com/haiyiyun/mongodb"
	"github.com/haiyiyun/webrouter"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	uploadConfFile := flag.String("config.plugins.upload", "../config/plugins/upload/upload.conf", "upload config file")
	var uploadConf service.Config
	config.Files(*uploadConfFile).Load(&uploadConf)

	uploadCache := cache.New(uploadConf.CacheDefaultExpiration.Duration, uploadConf.CacheCleanupInterval.Duration)
	uploadDB := mongodb.NewMongoPool("", uploadConf.MongoDatabaseName, 100, options.Client().ApplyURI(uploadConf.MongoDNS))
	webrouter.SetCloser(func() { uploadDB.Disconnect(context.TODO()) })

	uploadDB.M().InitCollection(schema.Upload)

	uploadService := service.NewService(&uploadConf, uploadCache, uploadDB)

	if uploadConf.WebRouter {
		//Init Begin
		uploadServiceUploadService := serviceUpload.NewService(uploadService)
		//Init End

		//Go Begin
		//Go End

		//Register Begin
		webrouter.Register(uploadConf.WebRouterRootPath+"upload/", uploadServiceUploadService)
		//Register End
	}

}
