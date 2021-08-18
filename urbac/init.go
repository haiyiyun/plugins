package urbac

import (
	"context"
	"flag"
	"net/http"

	"github.com/haiyiyun/plugins/urbac/database/schema"
	"github.com/haiyiyun/plugins/urbac/service"
	urbacServiceApplication "github.com/haiyiyun/plugins/urbac/service/application"
	urbacServiceAuth "github.com/haiyiyun/plugins/urbac/service/auth"
	urbacServiceProfile "github.com/haiyiyun/plugins/urbac/service/profile"
	urbacServiceRole "github.com/haiyiyun/plugins/urbac/service/role"
	urbacServiceToken "github.com/haiyiyun/plugins/urbac/service/token"
	urbacServiceUser "github.com/haiyiyun/plugins/urbac/service/user"

	"github.com/haiyiyun/cache"
	"github.com/haiyiyun/config"
	"github.com/haiyiyun/mongodb"
	"github.com/haiyiyun/utils/help"
	"github.com/haiyiyun/webrouter"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	urbacConfFile := flag.String("config.plugins.urbac", "../config/plugins/urbac/urbac.conf", "urbac config file")
	var urbacConf service.Config
	config.Files(*urbacConfFile).Load(&urbacConf)

	urbacCache := cache.New(urbacConf.CacheDefaultExpiration.Duration, urbacConf.CacheCleanupInterval.Duration)
	urbacDB := mongodb.NewMongoPool("", urbacConf.MongoDatabaseName, 100, options.Client().ApplyURI(urbacConf.MongoDNS))
	webrouter.SetCloser(func() { urbacDB.Disconnect(context.TODO()) })

	urbacDB.M().InitCollection(schema.User)
	urbacDB.M().InitCollection(schema.Role)
	urbacDB.M().InitCollection(schema.Application)
	urbacDB.M().InitCollection(schema.Token)

	urbacService := service.NewService(&urbacConf, urbacCache, urbacDB)

	if urbacConf.WebRouter {
		//Init Begin
		urbacServiceUserService := urbacServiceUser.NewService(urbacService)
		urbacServiceTokenService := urbacServiceToken.NewService(urbacService)
		urbacServiceRoleService := urbacServiceRole.NewService(urbacService)
		urbacServiceAuthService := urbacServiceAuth.NewService(urbacService)
		urbacServiceProfileService := urbacServiceProfile.NewService(urbacService)
		urbacServiceApplicationService := urbacServiceApplication.NewService(urbacService)
		//Init End

		//Go Begin
		//Go End

		//Register Begin
		webrouter.Register(urbacConf.WebRouterRootPath+"token/", urbacServiceTokenService)
		webrouter.Register(urbacConf.WebRouterRootPath+"user/", urbacServiceUserService)
		webrouter.Register(urbacConf.WebRouterRootPath+"role/", urbacServiceRoleService)
		webrouter.Register(urbacConf.WebRouterRootPath+"auth/", urbacServiceAuthService)
		webrouter.Register(urbacConf.WebRouterRootPath+"profile/", urbacServiceProfileService)
		webrouter.Register(urbacConf.WebRouterRootPath+"application/", urbacServiceApplicationService)
		//Register End
	}

	webrouter.Injector("urbac", "", 996, func(rw http.ResponseWriter, r *http.Request) (abort bool) {
		reqPath := r.URL.Path
		checkLogin := true
		if checkMethods, found := urbacConf.IgnoreCheckLoginPath[reqPath]; found {
			if len(checkMethods) == 0 || help.NewSlice(checkMethods).CheckPartItem(r.Method, "") {
				checkLogin = false
			}
		}

		if checkLogin {
			if u, found := urbacService.GetUserInfo(r); found {
				if allow := urbacService.CheckRight(reqPath, r.Method, u.ID); !allow {
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
