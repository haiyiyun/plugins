package cors

import (
	"net/http"
	"sync"

	"github.com/haiyiyun/webrouter"
)

var mut sync.Mutex

func init() {
	webrouter.Injector("cors", "", 999, func(rw http.ResponseWriter, r *http.Request) (abort bool) {
		mut.Lock()
		origin := r.Header.Get("Origin")
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token, X-Token, X-User-Id")
		rw.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST, GET, DELETE, PUT")
		rw.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
		mut.Unlock()

		if r.Method == "OPTIONS" {
			abort = true
			rw.WriteHeader(http.StatusNoContent)
		}

		return
	})
}
