package manage

import (
	"context"
	"flag"

	"github.com/haiyiyun/plugins/user_profile/database/schema"
	"github.com/haiyiyun/plugins/user_profile/service/base"
	"github.com/haiyiyun/plugins/user_profile/service/manage"
	manageService1 "github.com/haiyiyun/plugins/user_profile/service/manage/service1"

	"github.com/haiyiyun/cache"
	"github.com/haiyiyun/config"
	"github.com/haiyiyun/mongodb"
	"github.com/haiyiyun/webrouter"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	baseConfFile := flag.String("config.plugins.user_profile.manage.base", "../config/plugins/user_profile/base.conf", "base config file")
	var baseConf base.Config
	config.Files(*baseConfFile).Load(&baseConf)

	baseCache := cache.New(baseConf.CacheDefaultExpiration.Duration, baseConf.CacheCleanupInterval.Duration)
	baseDB := mongodb.NewMongoPool("", baseConf.MongoDatabaseName, 100, options.Client().ApplyURI(baseConf.MongoDNS))
	webrouter.SetCloser(func() { baseDB.Disconnect(context.TODO()) })

	baseDB.M().InitCollection(schema.Collection1)

	baseService := base.NewService(&baseConf, baseCache, baseDB)

	manageConfFile := flag.String("config.user_profile.manage", "../config/plugins/user_profile/manage.conf", "manage config file")
	var manageConf manage.Config
	config.Files(*manageConfFile).Load(&manageConf)

	if manageConf.WebRouter {
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
