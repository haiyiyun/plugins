package manage

import (
	"context"
	"flag"
	"os"

	"github.com/haiyiyun/plugins/content/database/schema"
	"github.com/haiyiyun/plugins/content/service/base"
	"github.com/haiyiyun/plugins/content/service/manage"

	"github.com/haiyiyun/cache"
	"github.com/haiyiyun/config"
	"github.com/haiyiyun/mongodb"
	"github.com/haiyiyun/webrouter"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	manageConfFile := flag.String("config.plugins.content.manage", "../config/plugins/content/manage.conf", "manage config file")
	var manageConf manage.Config
	config.Files(*manageConfFile).Load(&manageConf)

	if manageConf.WebRouter {
		baseConfFile := flag.String("config.plugins.content.manage.base", "../config/plugins/content/base.conf", "base config file")
		var baseConf base.Config
		config.Files(*baseConfFile).Load(&baseConf)

		os.Setenv("HYY_CACHE_TYPE", baseConf.CacheType)
		os.Setenv("HYY_CACHE_URL", baseConf.CacheUrl)
		os.Setenv("HYY_SHARD_COUNT", baseConf.CacheShardCount)
		os.Setenv("HYY_STRICT_TYPE_CHECK", baseConf.CacheUStrictTypeCheck)

		baseCache := cache.New(baseConf.CacheDefaultExpiration.Duration, baseConf.CacheCleanupInterval.Duration)
		baseDB := mongodb.NewMongoPool("", baseConf.MongoDatabaseName, 100, options.Client().ApplyURI(baseConf.MongoDNS))
		webrouter.SetCloser(func() { baseDB.Disconnect(context.TODO()) })

		baseDB.M().InitCollection(schema.Collection1)

		baseService := base.NewService(&baseConf, baseCache, baseDB)

		manageConf.Config = baseConf
		manageService := manage.NewService(&manageConf, baseService)

		//Init Begin
		// manageService1Service := manageService1.NewService(manageService)
		//Init End

		//Go Begin
		//Go End

		//Register Begin
		// webrouter.Register(manageConf.WebRouterRootPath+"", manageService1Service)
		//Register End
	}
}
