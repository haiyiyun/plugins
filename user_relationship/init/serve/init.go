package serve

import (
	"context"
	"flag"

	"github.com/haiyiyun/plugins/user_relationship/database/schema"
	"github.com/haiyiyun/plugins/user_relationship/service/base"
	userVisitor "github.com/haiyiyun/plugins/user_relationship/service/serve/visitor"

	"github.com/haiyiyun/cache"
	"github.com/haiyiyun/config"
	"github.com/haiyiyun/mongodb"
	"github.com/haiyiyun/plugins/user_relationship/service/serve"
	"github.com/haiyiyun/webrouter"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	baseConfFile := flag.String("config.plugins.user_relationship.serve.base", "../config/plugins/user_relationship/base.conf", "base config file")
	var baseConf base.Config
	config.Files(*baseConfFile).Load(&baseConf)

	baseCache := cache.New(baseConf.CacheDefaultExpiration.Duration, baseConf.CacheCleanupInterval.Duration)
	baseDB := mongodb.NewMongoPool("", baseConf.MongoDatabaseName, 100, options.Client().ApplyURI(baseConf.MongoDNS))
	webrouter.SetCloser(func() { baseDB.Disconnect(context.TODO()) })

	baseDB.M().InitCollection(schema.Contacts)
	baseDB.M().InitCollection(schema.ContactsApply)
	baseDB.M().InitCollection(schema.ContactsBlacklist)
	baseDB.M().InitCollection(schema.Group)
	baseDB.M().InitCollection(schema.GroupApply)
	baseDB.M().InitCollection(schema.GroupBlacklist)
	baseDB.M().InitCollection(schema.Visitor)

	baseService := base.NewService(&baseConf, baseCache, baseDB)

	serveConfFile := flag.String("config.user_relationship.serve", "../config/plugins/user_relationship/serve.conf", "serve config file")
	var serveConf serve.Config
	config.Files(*serveConfFile).Load(&serveConf)

	if serveConf.WebRouter {
		serveConf.Config = baseConf
		serveService := serve.NewService(&serveConf, baseService)

		//Init Begin
		userVisitorService := userVisitor.NewService(serveService)
		//Init End

		//Go Begin
		//Go End

		//Register Begin
		webrouter.Register(serveConf.WebRouterRootPath+"visitor/", userVisitorService)
		//Register End
	}
}
