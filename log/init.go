package log

import (
	"context"
	"flag"
	"net/http"

	"github.com/haiyiyun/plugins/log/database/schema"
	"github.com/haiyiyun/plugins/log/service"
	serviceLog "github.com/haiyiyun/plugins/log/service/log"

	"github.com/haiyiyun/cache"
	"github.com/haiyiyun/config"
	"github.com/haiyiyun/mongodb"
	"github.com/haiyiyun/webrouter"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	logConfFile := flag.String("config.plugins.log", "../config/plugins/log/log.conf", "log config file")
	var logConf service.Config
	config.Files(*logConfFile).Load(&logConf)

	if logConf.Log {
		logCache := cache.New(logConf.CacheDefaultExpiration.Duration, logConf.CacheCleanupInterval.Duration)
		logDB := mongodb.NewMongoPool("", logConf.MongoDatabaseName, 100, options.Client().ApplyURI(logConf.MongoDNS))
		webrouter.SetCloser(func() { logDB.Disconnect(context.TODO()) })

		logDB.M().InitCollection(schema.Log)

		logService := service.NewService(&logConf, logCache, logDB)

		if logConf.WebRouter {
			//Init Begin
			logServiceLogService := serviceLog.NewService(logService)
			//Init End

			//Go Begin
			//Go End

			//Register Begin
			webrouter.Register(logConf.WebRouterRootPath+"log/", logServiceLogService)
			//Register End
		}

		webrouter.Injector("loglogin", "", 997, func(rw http.ResponseWriter, r *http.Request) (abort bool) {
			if logIDHex := logService.LogRequestLogin(r); logIDHex != "" {
				lrw := &service.ResponseWriter{
					ResponseWriter: rw,
				}

				lrw.SetID(logIDHex)
				webrouter.ResponseWriter(lrw)
			}

			return
		})

		webrouter.Injector("logauth", "loglogin", 997, func(rw http.ResponseWriter, r *http.Request) (abort bool) {
			if logIDHex := logService.LogRequestAuth(r); logIDHex != "" {
				lrw := &service.ResponseWriter{
					ResponseWriter: rw,
				}

				lrw.SetID(logIDHex)
				webrouter.ResponseWriter(lrw)
			}

			return
		})

		webrouter.Injector("logoperate", "logauth", 997, func(rw http.ResponseWriter, r *http.Request) (abort bool) {
			if logIDHex := logService.LogRequestOperate(r); logIDHex != "" {
				lrw := &service.ResponseWriter{
					ResponseWriter: rw,
				}

				lrw.SetID(logIDHex)
				webrouter.ResponseWriter(lrw)
			}

			return
		})

		webrouter.Releasor("loglogin", "logauth", 1, func(rw http.ResponseWriter, r *http.Request) (abort bool) {
			logService.LogResponseLogin(rw, r)

			return
		})

		webrouter.Releasor("logauth", "logoperate", 1, func(rw http.ResponseWriter, r *http.Request) (abort bool) {
			logService.LogResponseAuth(rw, r)

			return
		})

		webrouter.Releasor("logoperate", "urbac", 1, func(rw http.ResponseWriter, r *http.Request) (abort bool) {
			logService.LogResponseOperate(rw, r)

			return
		})
	}
}
