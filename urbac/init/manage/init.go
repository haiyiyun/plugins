package urbac

import (
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/haiyiyun/plugins/urbac/database/schema"
	"github.com/haiyiyun/plugins/urbac/service/base"
	"github.com/haiyiyun/plugins/urbac/service/manage"
	manageApplication "github.com/haiyiyun/plugins/urbac/service/manage/application"
	manageAuth "github.com/haiyiyun/plugins/urbac/service/manage/auth"
	manageProfile "github.com/haiyiyun/plugins/urbac/service/manage/profile"
	manageRole "github.com/haiyiyun/plugins/urbac/service/manage/role"
	manageToken "github.com/haiyiyun/plugins/urbac/service/manage/token"
	manageUser "github.com/haiyiyun/plugins/urbac/service/manage/user"

	"github.com/haiyiyun/cache"
	"github.com/haiyiyun/config"
	"github.com/haiyiyun/mongodb"
	"github.com/haiyiyun/utils/help"
	"github.com/haiyiyun/webrouter"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	baseConfFile := flag.String("config.plugins.urbac.manage.base", "../config/plugins/urbac/base.conf", "base config file")
	var baseConf base.Config
	config.Files(*baseConfFile).Load(&baseConf)

	if baseConf.URBAC {
		os.Setenv("HYY_CACHE_TYPE", baseConf.CacheType)
		os.Setenv("HYY_CACHE_URL", baseConf.CacheUrl)
		os.Setenv("HYY_SHARD_COUNT", baseConf.CacheShardCount)
		os.Setenv("HYY_STRICT_TYPE_CHECK", baseConf.CacheUStrictTypeCheck)

		baseCache := cache.New(baseConf.CacheDefaultExpiration.Duration, baseConf.CacheCleanupInterval.Duration)
		baseDB := mongodb.NewMongoPool("", baseConf.MongoDatabaseName, 100, options.Client().ApplyURI(baseConf.MongoDNS))
		webrouter.SetCloser(func() { baseDB.Disconnect(context.TODO()) })

		baseDB.M().InitCollection(schema.User)
		baseDB.M().InitCollection(schema.Role)
		baseDB.M().InitCollection(schema.Application)
		baseDB.M().InitCollection(schema.Token)

		baseService := base.NewService(&baseConf, baseCache, baseDB)

		webrouter.Injector("urbac", "", 996, func(rw http.ResponseWriter, r *http.Request) (abort bool) {
			reqPath := r.URL.Path
			checkLogin := true
			if checkMethods, found := baseConf.IgnoreCheckLoginPath[reqPath]; found {
				if len(checkMethods) == 0 || help.NewSlice(checkMethods).CheckPartItem(r.Method, "") {
					checkLogin = false
				}
			}

			if checkLogin {
				if u, found := baseService.GetUserInfo(r); found {
					if allow := baseService.CheckRight(reqPath, r.Method, u.ID); !allow {
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

		manageConfFile := flag.String("config.plugins.urbac.manage", "../config/plugins/urbac/manage.conf", "manage config file")
		var manageConf manage.Config
		config.Files(*manageConfFile).Load(&manageConf)

		if manageConf.WebRouter {
			manageConf.Config = baseConf
			manageService := manage.NewService(&manageConf, baseService)

			//Init Begin
			manageTokenService := manageToken.NewService(manageService)
			manageUserService := manageUser.NewService(manageService)
			manageRoleService := manageRole.NewService(manageService)
			manageAuthService := manageAuth.NewService(manageService)
			manageProfileService := manageProfile.NewService(manageService)
			manageApplicationService := manageApplication.NewService(manageService)
			//Init End

			//Go Begin
			//Go End

			//Register Begin
			webrouter.Register(manageConf.WebRouterRootPath+"token/", manageTokenService)
			webrouter.Register(manageConf.WebRouterRootPath+"user/", manageUserService)
			webrouter.Register(manageConf.WebRouterRootPath+"role/", manageRoleService)
			webrouter.Register(manageConf.WebRouterRootPath+"auth/", manageAuthService)
			webrouter.Register(manageConf.WebRouterRootPath+"profile/", manageProfileService)
			webrouter.Register(manageConf.WebRouterRootPath+"application/", manageApplicationService)
			//Register End
		}
	}
}
