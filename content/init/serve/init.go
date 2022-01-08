package serve

import (
	"context"
	"flag"

	"github.com/haiyiyun/plugins/content/database/schema"
	"github.com/haiyiyun/plugins/content/service/base"
	contentCategory "github.com/haiyiyun/plugins/content/service/serve/category"
	contentContent "github.com/haiyiyun/plugins/content/service/serve/content"
	contentDiscuss "github.com/haiyiyun/plugins/content/service/serve/discuss"
	contentFollow "github.com/haiyiyun/plugins/content/service/serve/follow"
	contentSubject "github.com/haiyiyun/plugins/content/service/serve/subject"

	"github.com/haiyiyun/cache"
	"github.com/haiyiyun/config"
	"github.com/haiyiyun/mongodb"
	"github.com/haiyiyun/plugins/content/service/serve"
	"github.com/haiyiyun/webrouter"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	serveConfFile := flag.String("config.content.serve", "../config/plugins/content/serve.conf", "serve config file")
	var serveConf serve.Config
	config.Files(*serveConfFile).Load(&serveConf)

	if serveConf.WebRouter {
		baseConfFile := flag.String("config.plugins.content.serve.base", "../config/plugins/content/base.conf", "base config file")
		var baseConf base.Config
		config.Files(*baseConfFile).Load(&baseConf)

		baseCache := cache.New(baseConf.CacheDefaultExpiration.Duration, baseConf.CacheCleanupInterval.Duration)
		baseDB := mongodb.NewMongoPool("", baseConf.MongoDatabaseName, 100, options.Client().ApplyURI(baseConf.MongoDNS))
		webrouter.SetCloser(func() { baseDB.Disconnect(context.TODO()) })

		baseDB.M().InitCollection(schema.Category)
		baseDB.M().InitCollection(schema.Subject)
		baseDB.M().InitCollection(schema.Content)
		baseDB.M().InitCollection(schema.Discuss)
		baseDB.M().InitCollection(schema.FollowRelationship)
		baseDB.M().InitCollection(schema.FollowContent)
		baseDB.M().InitCollection(schema.PublishIntervene)
		baseDB.M().InitCollection(schema.Favorites)
		baseDB.M().InitCollection(schema.KeywordBan)
		baseDB.M().InitCollection(schema.Message)
		baseDB.M().InitCollection(schema.GroupMessage)
		baseDB.M().InitCollection(schema.Share)

		baseService := base.NewService(&baseConf, baseCache, baseDB)

		serveConf.Config = baseConf
		serveService := serve.NewService(&serveConf, baseService)

		//Init Begin
		contentCategoryService := contentCategory.NewService(serveService)
		contentSubjectService := contentSubject.NewService(serveService)
		contentContentService := contentContent.NewService(serveService)
		contentDiscussService := contentDiscuss.NewService(serveService)
		contentFollowService := contentFollow.NewService(serveService)
		//Init End

		//Go Begin
		//Go End

		//Register Begin
		webrouter.Register(serveConf.WebRouterRootPath+"category/", contentCategoryService)
		webrouter.Register(serveConf.WebRouterRootPath+"subject/", contentSubjectService)
		webrouter.Register(serveConf.WebRouterRootPath+"content/", contentContentService)
		webrouter.Register(serveConf.WebRouterRootPath+"discuss/", contentDiscussService)
		webrouter.Register(serveConf.WebRouterRootPath+"follow/", contentFollowService)
		//Register End
	}
}
