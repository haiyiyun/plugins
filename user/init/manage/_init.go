package manage

import (
	"context"
	"flag"

	"github.com/haiyiyun/plugins/user/database/schema"
	"github.com/haiyiyun/plugins/user/service/base"
	"github.com/haiyiyun/plugins/user/service/manage"
	manageService1 "github.com/haiyiyun/plugins/user/service/manage/service1"

	"github.com/haiyiyun/cache"
	"github.com/haiyiyun/config"
	"github.com/haiyiyun/mongodb"
	"github.com/haiyiyun/webrouter"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	manageConfFile := flag.String("config.webrouter_plugin_template.manage", "../config/plugins/webrouter_plugin_template/manage.conf", "manage config file")
	var manageConf manage.Config
	config.Files(*manageConfFile).Load(&manageConf)

	if manageConf.WebRouter {
		baseConfFile := flag.String("config.plugins.webrouter_plugin_template.manage.base", "../config/plugins/webrouter_plugin_template/base.conf", "base config file")
		var baseConf base.Config
		config.Files(*baseConfFile).Load(&baseConf)

		baseCache := cache.New(baseConf.CacheDefaultExpiration.Duration, baseConf.CacheCleanupInterval.Duration)
		baseDB := mongodb.NewMongoPool("", baseConf.MongoDatabaseName, 100, options.Client().ApplyURI(baseConf.MongoDNS))
		webrouter.SetCloser(func() { baseDB.Disconnect(context.TODO()) })

		baseDB.M().InitCollection(schema.Collection1)

		baseService := base.NewService(&baseConf, baseCache, baseDB)

		manageConf.Config = baseConf
		manageService := manage.NewService(&manageConf, baseService)

		//Init Begin
		manageService1Service := manageService1.NewService(manageService)
		//Init End

		//Go Begin
		//Go End

		//Register Begin
		webrouter.Register(manageConf.WebRouterRootPath+"", manageService1Service)
		//Register End
	}
}
