package upload

import (
	"context"
	"flag"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/haiyiyun/plugins/upload/service/aliyun"   // 导入阿里云存储
	_ "github.com/haiyiyun/plugins/upload/service/local"    // 导入本地存储
	_ "github.com/haiyiyun/plugins/upload/service/qiniuyun" // 导入七牛云存储
	_ "github.com/haiyiyun/plugins/upload/service/tencent"  // 导入腾讯云存储

	"github.com/haiyiyun/plugins/upload/database/schema"
	"github.com/haiyiyun/plugins/upload/service/base"
	"github.com/haiyiyun/plugins/upload/service/serve"
	serveUpload "github.com/haiyiyun/plugins/upload/service/serve/upload"

	"github.com/haiyiyun/cache"
	"github.com/haiyiyun/config"
	"github.com/haiyiyun/log"
	"github.com/haiyiyun/mongodb"
	"github.com/haiyiyun/webrouter"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	serveConfFile := flag.String("config.plugins.upload.serve", "../config/plugins/upload/serve.conf", "serve config file")
	var serveConf serve.Config
	config.Files(*serveConfFile).Load(&serveConf)

	if serveConf.WebRouter {
		baseConfFile := flag.String("config.plugins.upload.serve.base", "../config/plugins/upload/base.conf", "base config file")
		var baseConf base.Config
		config.Files(*baseConfFile).Load(&baseConf)

		uploadDir := filepath.Clean(baseConf.UploadDirectory)
		if _, err := os.Stat(uploadDir); err != nil {
			log.Fatal("upload directory must exist and only manually create")
		}

		// 验证必要的配置项
		if baseConf.StorageType == "" {
			log.Fatal("storage_type must be configured")
		}

		os.Setenv("HYY_CACHE_TYPE", baseConf.CacheType)
		os.Setenv("HYY_CACHE_URL", baseConf.CacheUrl)
		os.Setenv("HYY_SHARD_COUNT", baseConf.CacheShardCount)
		os.Setenv("HYY_STRICT_TYPE_CHECK", baseConf.CacheUStrictTypeCheck)

		baseCache := cache.New(baseConf.CacheDefaultExpiration.Duration, baseConf.CacheCleanupInterval.Duration)
		baseDB := mongodb.NewMongoPool("", baseConf.MongoDatabaseName, 100, options.Client().ApplyURI(baseConf.MongoDNS))
		webrouter.SetCloser(func() { baseDB.Disconnect(context.TODO()) })

		baseDB.M().InitCollection(schema.Upload)

		baseService := base.NewService(&baseConf, baseCache, baseDB)

		serveConf.Config = baseConf

		if serveConf.MaxUploadFileSize == 0 {
			//如果没设置则使用默认32M
			serveConf.MaxUploadFileSize = 32 << 20
		}

		if serveConf.BuildInFileServer && serveConf.AllowDownloadLocal {
			webrouter.Handle(serveConf.DownloadLocalUrlDirectory, http.StripPrefix(serveConf.DownloadLocalUrlDirectory, http.FileServer(http.Dir(baseConf.UploadDirectory))))
		}

		serveService := serve.NewService(&serveConf, baseService)

		//Init Begin
		serveUploadService := serveUpload.NewService(serveService)
		//Init End

		//Go Begin
		//Go End

		//Register Begin
		webrouter.Register(serveConf.WebRouterRootPath+"upload/", serveUploadService)
		//Register End
	}

}
