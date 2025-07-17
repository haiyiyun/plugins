package serve

import (
	"context"
	"flag"
	"os"

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
	serveConfFile := flag.String("config.plugins.user_relationship.serve", "../config/plugins/user_relationship/serve.conf", "serve config file")
	var serveConf serve.Config
	config.Files(*serveConfFile).Load(&serveConf)

	if serveConf.WebRouter {
		baseConfFile := flag.String("config.plugins.user_relationship.serve.base", "../config/plugins/user_relationship/base.conf", "base config file")
		var baseConf base.Config
		config.Files(*baseConfFile).Load(&baseConf)

		os.Setenv("HYY_CACHE_TYPE", baseConf.CacheType)
		os.Setenv("HYY_CACHE_URL", baseConf.CacheUrl)
		os.Setenv("HYY_SHARD_COUNT", baseConf.CacheShardCount)
		os.Setenv("HYY_STRICT_TYPE_CHECK", baseConf.CacheUStrictTypeCheck)

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
