package serve

import (
	"context"
	"flag"
	"os"

	"github.com/haiyiyun/plugins/cities/database/schema"
	"github.com/haiyiyun/plugins/cities/service/base"
	"github.com/haiyiyun/plugins/cities/service/serve"
	serveCities "github.com/haiyiyun/plugins/cities/service/serve/cities"

	"github.com/haiyiyun/cache"
	"github.com/haiyiyun/config"
	"github.com/haiyiyun/mongodb"
	"github.com/haiyiyun/webrouter"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	serveConfFile := flag.String("config.plugins.cities.serve", "../config/plugins/cities/serve.conf", "serve config file")
	var serveConf serve.Config
	config.Files(*serveConfFile).Load(&serveConf)

	if serveConf.WebRouter {
		baseConfFile := flag.String("config.plugins.cities.serve.base", "../config/plugins/cities/base.conf", "base config file")
		var baseConf base.Config
		config.Files(*baseConfFile).Load(&baseConf)

		os.Setenv("HYY_CACHE_TYPE", baseConf.CacheType)
		os.Setenv("HYY_CACHE_URL", baseConf.CacheUrl)

		baseCache := cache.New(baseConf.CacheDefaultExpiration.Duration, baseConf.CacheCleanupInterval.Duration)
		baseDB := mongodb.NewMongoPool("", baseConf.MongoDatabaseName, 100, options.Client().ApplyURI(baseConf.MongoDNS))
		webrouter.SetCloser(func() { baseDB.Disconnect(context.TODO()) })

		baseDB.M().InitCollection(schema.Province)
		baseDB.M().InitCollection(schema.City)
		baseDB.M().InitCollection(schema.Area)
		baseDB.M().InitCollection(schema.Street)
		baseDB.M().InitCollection(schema.Village)

		baseService := base.NewService(&baseConf, baseCache, baseDB)

		serveConf.Config = baseConf
		serveService := serve.NewService(&serveConf, baseService)

		//Init Begin
		serveCitiesService := serveCities.NewService(serveService)
		//Init End

		//Go Begin
		//Go End

		//Register Begin
		webrouter.Register(serveConf.WebRouterRootPath+"", serveCitiesService)
		//Register End
	}
}
