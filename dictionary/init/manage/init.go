package upload

import (
	"context"
	"flag"
	"os"

	"github.com/haiyiyun/plugins/dictionary/database/schema"
	"github.com/haiyiyun/plugins/dictionary/service/base"
	"github.com/haiyiyun/plugins/dictionary/service/manage"
	manageDictionary "github.com/haiyiyun/plugins/dictionary/service/manage/dictionary"

	"github.com/haiyiyun/cache"
	"github.com/haiyiyun/config"
	"github.com/haiyiyun/mongodb"
	"github.com/haiyiyun/webrouter"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	manageConfFile := flag.String("config.plugins.dictionary.manage", "../config/plugins/dictionary/manage.conf", "manage config file")
	var manageConf manage.Config
	config.Files(*manageConfFile).Load(&manageConf)

	if manageConf.WebRouter {
		baseConfFile := flag.String("config.plugins.dictionary.manage.base", "../config/plugins/dictionary/base.conf", "base config file")
		var baseConf base.Config
		config.Files(*baseConfFile).Load(&baseConf)

		os.Setenv("HYY_CACHE_TYPE", baseConf.CacheType)
		os.Setenv("HYY_CACHE_URL", baseConf.CacheUrl)

		baseCache := cache.New(baseConf.CacheDefaultExpiration.Duration, baseConf.CacheCleanupInterval.Duration)
		baseDB := mongodb.NewMongoPool("", baseConf.MongoDatabaseName, 100, options.Client().ApplyURI(baseConf.MongoDNS))
		webrouter.SetCloser(func() { baseDB.Disconnect(context.TODO()) })

		baseDB.M().InitCollection(schema.Dictionary)

		baseService := base.NewService(&baseConf, baseCache, baseDB)

		manageConf.Config = baseConf
		manageService := manage.NewService(&manageConf, baseService)

		//Init Begin
		manageDictionaryService := manageDictionary.NewService(manageService)
		//Init End

		//Go Begin
		//Go End

		//Register Begin
		webrouter.Register(manageConf.WebRouterRootPath+"", manageDictionaryService)
		//Register End
	}

}
