package serve

import (
	"context"
	"flag"
	"net/http"
	"os"

	"github.com/haiyiyun/plugins/log/database/schema"
	"github.com/haiyiyun/plugins/log/service/base"
	"github.com/haiyiyun/plugins/log/service/serve"

	"github.com/haiyiyun/cache"
	"github.com/haiyiyun/config"
	"github.com/haiyiyun/mongodb"
	"github.com/haiyiyun/webrouter"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	serveConfFile := flag.String("config.plugins.log.serve", "../config/plugins/log/serve.conf", "serve config file")
	var serveConf serve.Config
	config.Files(*serveConfFile).Load(&serveConf)

	if serveConf.Log {
		baseConfFile := flag.String("config.plugins.log.serve.base", "../config/plugins/log/base.conf", "base config file")
		var baseConf base.Config
		config.Files(*baseConfFile).Load(&baseConf)

		os.Setenv("HYY_CACHE_TYPE", baseConf.CacheType)
		os.Setenv("HYY_CACHE_URL", baseConf.CacheUrl)
		os.Setenv("HYY_SHARD_COUNT", baseConf.CacheShardCount)
		os.Setenv("HYY_STRICT_TYPE_CHECK", baseConf.CacheUStrictTypeCheck)

		baseCache := cache.New(baseConf.CacheDefaultExpiration.Duration, baseConf.CacheCleanupInterval.Duration)
		baseDB := mongodb.NewMongoPool("", baseConf.MongoDatabaseName, 100, options.Client().ApplyURI(baseConf.MongoDNS))
		webrouter.SetCloser(func() { baseDB.Disconnect(context.TODO()) })

		baseDB.M().InitCollection(schema.Log)

		baseService := base.NewService(&baseConf, baseCache, baseDB)

		webrouter.Injector("loglogin", "", 997, func(rw http.ResponseWriter, r *http.Request) (abort bool) {
			if logID := baseService.LogRequestLogin(r); logID != primitive.NilObjectID {
				if lrw, ok := rw.(*webrouter.ResponseWriter); ok {
					lrw.SetGetResData(true)
					lrw.SetData("log_id", logID)
				}
			}

			return
		})

		webrouter.Injector("logauth", "loglogin", 997, func(rw http.ResponseWriter, r *http.Request) (abort bool) {
			if logID := baseService.LogRequestAuth(r); logID != primitive.NilObjectID {
				if lrw, ok := rw.(*webrouter.ResponseWriter); ok {
					lrw.SetGetResData(true)
					lrw.SetData("log_id", logID)
				}
			}

			return
		})

		webrouter.Injector("logoperate", "logauth", 997, func(rw http.ResponseWriter, r *http.Request) (abort bool) {
			if logID := baseService.LogRequestOperate(r); logID != primitive.NilObjectID {
				if lrw, ok := rw.(*webrouter.ResponseWriter); ok {
					lrw.SetGetResData(true)
					lrw.SetData("log_id", logID)
				}
			}

			return
		})

		webrouter.Releasor("loglogin", "logauth", 1, func(rw http.ResponseWriter, r *http.Request) (abort bool) {
			baseService.LogResponseLogin(rw, r)
			return
		})

		webrouter.Releasor("logauth", "logoperate", 1, func(rw http.ResponseWriter, r *http.Request) (abort bool) {
			baseService.LogResponseAuth(rw, r)
			return
		})

		webrouter.Releasor("logoperate", "user", 1, func(rw http.ResponseWriter, r *http.Request) (abort bool) {
			baseService.LogResponseOperate(rw, r)
			return
		})
	}
}
