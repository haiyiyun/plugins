package serve

import (
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/haiyiyun/plugins/user/database/schema"
	"github.com/haiyiyun/plugins/user/service/base"
	userAuth "github.com/haiyiyun/plugins/user/service/serve/auth"
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

	if baseConf.CheckLogin {
		os.Setenv("HYY_CACHE_TYPE", baseConf.CacheType)
		os.Setenv("HYY_CACHE_URL", baseConf.CacheUrl)
		os.Setenv("HYY_SHARD_COUNT", baseConf.CacheShardCount)
		os.Setenv("HYY_STRICT_TYPE_CHECK", baseConf.CacheUStrictTypeCheck)

		baseCache := cache.New(baseConf.CacheDefaultExpiration.Duration, baseConf.CacheCleanupInterval.Duration)
		baseDB := mongodb.NewMongoPool("", baseConf.MongoDatabaseName, 100, options.Client().ApplyURI(baseConf.MongoDNS))
		webrouter.SetCloser(func() { baseDB.Disconnect(context.TODO()) })

		baseDB.M().InitCollection(schema.User)
		baseDB.M().InitCollection(schema.Token)

		baseService := base.NewService(&baseConf, baseCache, baseDB)

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

		serveConfFile := flag.String("config.plugins.user.serve", "../config/plugins/user/serve.conf", "serve config file")
		var serveConf serve.Config
		config.Files(*serveConfFile).Load(&serveConf)

		if serveConf.WebRouter {
			serveConf.Config = baseConf
			serveService := serve.NewService(&serveConf, baseService)
			if serveConf.EnableProfile {
				baseDB.M().InitCollection(schema.Profile)
			}

			//Init Begin
			userAuthService := userAuth.NewService(serveService)
			userUserService := userUser.NewService(serveService)
			//Init End

			//Go Begin
			//Go End

			//Register Begin
			webrouter.Register(serveConf.WebRouterRootPath+"auth/", userAuthService)
			webrouter.Register(serveConf.WebRouterRootPath+"user/", userUserService)
			//Register End
		}
	}
}
