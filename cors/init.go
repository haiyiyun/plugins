package cors

import (
	"flag"
	"net/http"
	"os"

	"github.com/haiyiyun/config"
	"github.com/haiyiyun/webrouter"
)

func init() {
	confFile := flag.String("config.plugins.cors", "../config/plugins/cors/cors.conf", "cors config file")

	var conf Config
	if _, err := os.Stat(*confFile); err == nil {
		config.Files(*confFile).Load(&conf)
	}

	webrouter.Injector("cors", "", 99999, func(rw http.ResponseWriter, r *http.Request) (abort bool) {
		origin := conf.AccessControlAllowOrigin
		if origin == "" {
			origin = r.Header.Get("Origin")
		}
		rw.Header().Set("Access-Control-Allow-Origin", origin)

		accessControlAllowHeaders := conf.AccessControlAllowHeaders
		if accessControlAllowHeaders == "" {
			accessControlAllowHeaders = "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token, X-Token, X-User-Id"
		}
		rw.Header().Set("Access-Control-Allow-Headers", accessControlAllowHeaders)

		accessControlAllowMethods := conf.AccessControlAllowMethods
		if accessControlAllowMethods == "" {
			accessControlAllowMethods = "OPTIONS, POST, GET, DELETE, PUT"
		}
		rw.Header().Set("Access-Control-Allow-Methods", accessControlAllowMethods)

		accessControlExposeHeaders := conf.AccessControlExposeHeaders
		if accessControlExposeHeaders == "" {
			accessControlExposeHeaders = "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type"
		}
		rw.Header().Set("Access-Control-Expose-Headers", accessControlExposeHeaders)

		accessControlAllowCredentials := conf.AccessControlAllowCredentials
		if accessControlAllowCredentials == "" {
			accessControlAllowCredentials = "true"
		}
		rw.Header().Set("Access-Control-Allow-Credentials", accessControlAllowCredentials)

		if r.Method == "OPTIONS" {
			abort = true
			rw.WriteHeader(http.StatusNoContent)
		}

		return
	})
}
