package upload

import (
	"context"
	"flag"
	"os"

	"github.com/haiyiyun/plugins/dictionary/database/schema"
	"github.com/haiyiyun/plugins/dictionary/service/base"
	"github.com/haiyiyun/plugins/dictionary/service/serve"
	serveDictionary "github.com/haiyiyun/plugins/dictionary/service/serve/dictionary"

	"github.com/haiyiyun/cache"
	"github.com/haiyiyun/config"
	"github.com/haiyiyun/mongodb"
	"github.com/haiyiyun/webrouter"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	serveConfFile := flag.String("config.plugins.dictionary.serve", "../config/plugins/dictionary/serve.conf", "serve config file")
	var serveConf serve.Config
	config.Files(*serveConfFile).Load(&serveConf)

	if serveConf.WebRouter {
		baseConfFile := flag.String("config.plugins.dictionary.serve.base", "../config/plugins/dictionary/base.conf", "base config file")
		var baseConf base.Config
		config.Files(*baseConfFile).Load(&baseConf)

		os.Setenv("HYY_CACHE_TYPE", baseConf.CacheType)
		os.Setenv("HYY_CACHE_URL", baseConf.CacheUrl)

		baseCache := cache.New(baseConf.CacheDefaultExpiration.Duration, baseConf.CacheCleanupInterval.Duration)
		baseDB := mongodb.NewMongoPool("", baseConf.MongoDatabaseName, 100, options.Client().ApplyURI(baseConf.MongoDNS))
		webrouter.SetCloser(func() { baseDB.Disconnect(context.TODO()) })

		baseDB.M().InitCollection(schema.Dictionary)

		baseService := base.NewService(&baseConf, baseCache, baseDB)

		serveConf.Config = baseConf
		serveService := serve.NewService(&serveConf, baseService)

		//Init Begin
		serveDictionaryService := serveDictionary.NewService(serveService)
		//Init End

		//Go Begin
		//Go End

		//Register Begin
		webrouter.Register(serveConf.WebRouterRootPath+"", serveDictionaryService)
		//Register End
	}

}
