package serve

import (
	"context"
	"flag"

	"github.com/haiyiyun/plugins/user_profile/database/schema"
	"github.com/haiyiyun/plugins/user_profile/service/base"
	userProfile "github.com/haiyiyun/plugins/user_profile/service/serve/profile"
	userPublic "github.com/haiyiyun/plugins/user_profile/service/serve/public"

	"github.com/haiyiyun/cache"
	"github.com/haiyiyun/config"
	"github.com/haiyiyun/mongodb"
	"github.com/haiyiyun/plugins/user_profile/service/serve"
	"github.com/haiyiyun/webrouter"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	baseConfFile := flag.String("config.plugins.user_profile.serve.base", "../config/plugins/user_profile/base.conf", "base config file")
	var baseConf base.Config
	config.Files(*baseConfFile).Load(&baseConf)

	baseCache := cache.New(baseConf.CacheDefaultExpiration.Duration, baseConf.CacheCleanupInterval.Duration)
	baseDB := mongodb.NewMongoPool("", baseConf.MongoDatabaseName, 100, options.Client().ApplyURI(baseConf.MongoDNS))
	webrouter.SetCloser(func() { baseDB.Disconnect(context.TODO()) })

	baseDB.M().InitCollection(schema.Audit)
	baseDB.M().InitCollection(schema.Profile)
	baseDB.M().InitCollection(schema.ProfileAudit)

	baseService := base.NewService(&baseConf, baseCache, baseDB)

	serveConfFile := flag.String("config.user_profile.serve", "../config/plugins/user_profile/serve.conf", "serve config file")
	var serveConf serve.Config
	config.Files(*serveConfFile).Load(&serveConf)

	if serveConf.WebRouter {
		serveConf.Config = baseConf
		serveService := serve.NewService(&serveConf, baseService)

		//Init Begin
		userProfileService := userProfile.NewService(serveService)
		userPublicService := userPublic.NewService(serveService)
		//Init End

		//Go Begin
		//Go End

		//Register Begin
		webrouter.Register(serveConf.WebRouterRootPath+"profile/", userProfileService)
		webrouter.Register(serveConf.WebRouterRootPath+"public/", userPublicService)
		//Register End
	}
}
