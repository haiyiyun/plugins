package serve

import (
	"context"
	"flag"
	"net/http"

	"github.com/haiyiyun/plugins/user/database/schema"
	"github.com/haiyiyun/plugins/user/service/base"
	userAuth "github.com/haiyiyun/plugins/user/service/serve/auth"
	userProfile "github.com/haiyiyun/plugins/user/service/serve/profile"
	userUser "github.com/haiyiyun/plugins/user/service/serve/user"

	"github.com/haiyiyun/cache"
	"github.com/haiyiyun/config"
	"github.com/haiyiyun/mongodb"
	"github.com/haiyiyun/plugins/user/service/serve"
	"github.com/haiyiyun/utils/help"
	"github.com/haiyiyun/webrouter"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	baseConfFile := flag.String("config.plugins.user.serve.base", "../config/plugins/user/base.conf", "base config file")
	var baseConf base.Config
	config.Files(*baseConfFile).Load(&baseConf)

	baseCache := cache.New(baseConf.CacheDefaultExpiration.Duration, baseConf.CacheCleanupInterval.Duration)
	baseDB := mongodb.NewMongoPool("", baseConf.MongoDatabaseName, 100, options.Client().ApplyURI(baseConf.MongoDNS))
	webrouter.SetCloser(func() { baseDB.Disconnect(context.TODO()) })

	baseDB.M().InitCollection(schema.User)
	baseDB.M().InitCollection(schema.Token)
	baseDB.M().InitCollection(schema.Profile)
	baseDB.M().InitCollection(schema.Audit)
	baseDB.M().InitCollection(schema.ProfileAudit)
	baseDB.M().InitCollection(schema.Contacts)
	baseDB.M().InitCollection(schema.ContactsApply)
	baseDB.M().InitCollection(schema.ContactsBlacklist)
	baseDB.M().InitCollection(schema.Group)
	baseDB.M().InitCollection(schema.GroupApply)
	baseDB.M().InitCollection(schema.GroupBlacklist)
	baseDB.M().InitCollection(schema.Visitor)

	baseService := base.NewService(&baseConf, baseCache, baseDB)

	if baseConf.CheckLogin {
		webrouter.Injector("user", "", 996, func(rw http.ResponseWriter, r *http.Request) (abort bool) {
			reqPath := r.URL.Path
			checkLogin := true
			if checkMethods, found := baseConf.IgnoreCheckLoginPath[reqPath]; found {
				if len(checkMethods) == 0 || help.NewSlice(checkMethods).CheckPartItem(r.Method, "") {
					checkLogin = false
				}
			}

			if checkLogin {
				if _, found := baseService.GetUserInfo(r); !found {
					rw.WriteHeader(http.StatusUnauthorized)
					return true
				}
			}

			return
		})
	}

	serveConfFile := flag.String("config.user.serve", "../config/plugins/user/serve.conf", "serve config file")
	var serveConf serve.Config
	config.Files(*serveConfFile).Load(&serveConf)

	if serveConf.WebRouter {
		serveConf.Config = baseConf
		serveService := serve.NewService(&serveConf, baseService)

		//Init Begin
		userAuthService := userAuth.NewService(serveService)
		userUserService := userUser.NewService(serveService)
		userProfileService := userProfile.NewService(serveService)
		//Init End

		//Go Begin
		//Go End

		//Register Begin
		webrouter.Register(serveConf.WebRouterRootPath+"auth/", userAuthService)
		webrouter.Register(serveConf.WebRouterRootPath+"user/", userUserService)
		webrouter.Register(serveConf.WebRouterRootPath+"profile/", userProfileService)
		//Register End
	}
}
