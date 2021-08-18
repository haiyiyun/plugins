package urbac

import (
	"context"
	"flag"
	"net/http"

	"github.com/haiyiyun/plugins/urbac/database/schema"
	"github.com/haiyiyun/plugins/urbac/service"
	userRBACServiceApplication "github.com/haiyiyun/plugins/urbac/service/application"
	userRBACServiceAuth "github.com/haiyiyun/plugins/urbac/service/auth"
	userRBACServiceProfile "github.com/haiyiyun/plugins/urbac/service/profile"
	userRBACServiceRole "github.com/haiyiyun/plugins/urbac/service/role"
	userRBACServiceToken "github.com/haiyiyun/plugins/urbac/service/token"
	userRBACServiceUser "github.com/haiyiyun/plugins/urbac/service/user"

	"github.com/haiyiyun/cache"
	"github.com/haiyiyun/config"
	"github.com/haiyiyun/mongodb"
	"github.com/haiyiyun/utils/help"
	"github.com/haiyiyun/webrouter"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	userRBACConfFile := flag.String("config.plugins.urbac", "../config/plugins/urbac/urbac.conf", "urbac config file")
	var userRBACConf service.Config
	config.Files(*userRBACConfFile).Load(&userRBACConf)

	userRBACCache := cache.New(userRBACConf.CacheDefaultExpiration.Duration, userRBACConf.CacheCleanupInterval.Duration)
	userRBACDB := mongodb.NewMongoPool("", userRBACConf.MongoDatabaseName, 100, options.Client().ApplyURI(userRBACConf.MongoDNS))
	webrouter.SetCloser(func() { userRBACDB.Disconnect(context.TODO()) })

	userRBACDB.M().InitCollection(schema.User)
	userRBACDB.M().InitCollection(schema.Role)
	userRBACDB.M().InitCollection(schema.Application)
	userRBACDB.M().InitCollection(schema.Token)

	userRBACService := service.NewService(&userRBACConf, userRBACCache, userRBACDB)

	if userRBACConf.WebRouter {
		//Init Begin
		userRBACServiceUserService := userRBACServiceUser.NewService(userRBACService)
		userRBACServiceTokenService := userRBACServiceToken.NewService(userRBACService)
		userRBACServiceRoleService := userRBACServiceRole.NewService(userRBACService)
		userRBACServiceAuthService := userRBACServiceAuth.NewService(userRBACService)
		userRBACServiceProfileService := userRBACServiceProfile.NewService(userRBACService)
		userRBACServiceApplicationService := userRBACServiceApplication.NewService(userRBACService)
		//Init End

		//Go Begin
		//Go End

		//Register Begin
		webrouter.Register(userRBACConf.WebRouterRootPath+"token/", userRBACServiceTokenService)
		webrouter.Register(userRBACConf.WebRouterRootPath+"user/", userRBACServiceUserService)
		webrouter.Register(userRBACConf.WebRouterRootPath+"role/", userRBACServiceRoleService)
		webrouter.Register(userRBACConf.WebRouterRootPath+"auth/", userRBACServiceAuthService)
		webrouter.Register(userRBACConf.WebRouterRootPath+"profile/", userRBACServiceProfileService)
		webrouter.Register(userRBACConf.WebRouterRootPath+"application/", userRBACServiceApplicationService)
		//Register End
	}

	webrouter.Injector("urbac", "", 996, func(rw http.ResponseWriter, r *http.Request) (abort bool) {
		reqPath := r.URL.Path
		checkLogin := true
		if checkMethods, found := userRBACConf.IgnoreCheckLoginPath[reqPath]; found {
			if len(checkMethods) == 0 || help.NewSlice(checkMethods).CheckPartItem(r.Method, "") {
				checkLogin = false
			}
		}

		if checkLogin {
			if u, found := userRBACService.GetUserInfo(r); found {
				if allow := userRBACService.CheckRight(reqPath, r.Method, u.ID); !allow {
					rw.WriteHeader(http.StatusForbidden)
					return true
				}
			} else {
				rw.WriteHeader(http.StatusUnauthorized)
				return true
			}
		}

		return
	})
}
