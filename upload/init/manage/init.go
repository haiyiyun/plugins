package upload

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/haiyiyun/plugins/upload/database/schema"
	"github.com/haiyiyun/plugins/upload/service/base"
	"github.com/haiyiyun/plugins/upload/service/manage"
	manageUpload "github.com/haiyiyun/plugins/upload/service/manage/upload"

	"github.com/haiyiyun/cache"
	"github.com/haiyiyun/config"
	"github.com/haiyiyun/mongodb"
	"github.com/haiyiyun/webrouter"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	manageConfFile := flag.String("config.plugins.upload.manage", "../config/plugins/upload/manage.conf", "manage config file")
	var manageConf manage.Config
	config.Files(*manageConfFile).Load(&manageConf)

	if manageConf.WebRouter {
		baseConfFile := flag.String("config.plugins.upload.manage.base", "../config/plugins/upload/base.conf", "base config file")
		var baseConf base.Config
		config.Files(*baseConfFile).Load(&baseConf)

		uploadDir := filepath.Clean(baseConf.UploadDirectory)
		if _, err := os.Stat(uploadDir); err != nil {
			log.Fatal("upload directory must exist and only manually create")
		}

		baseCache := cache.New(baseConf.CacheDefaultExpiration.Duration, baseConf.CacheCleanupInterval.Duration)
		baseDB := mongodb.NewMongoPool("", baseConf.MongoDatabaseName, 100, options.Client().ApplyURI(baseConf.MongoDNS))
		webrouter.SetCloser(func() { baseDB.Disconnect(context.TODO()) })

		baseDB.M().InitCollection(schema.Upload)

		baseService := base.NewService(&baseConf, baseCache, baseDB)

		manageConf.Config = baseConf
		if manageConf.BuildInFileServer && manageConf.AllowDownloadLocal {
			webrouter.Handle(manageConf.DownloadLocalUrlDirectory, http.StripPrefix(manageConf.DownloadLocalUrlDirectory, http.FileServer(http.Dir(baseConf.UploadDirectory))))
		}

		manageService := manage.NewService(&manageConf, baseService)

		//Init Begin
		manageUploadService := manageUpload.NewService(manageService)
		//Init End

		//Go Begin
		//Go End

		//Register Begin
		webrouter.Register(manageConf.WebRouterRootPath+"upload/", manageUploadService)
		//Register End
	}

}
