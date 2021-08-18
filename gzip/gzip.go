package gzip

import (
	"compress/gzip"
	"net/http"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/webrouter"
)

func init() {
	webrouter.Injector("gzip", "", 998, func(rw http.ResponseWriter, r *http.Request) (abort bool) {
		accept := r.Header.Get("Content-Encoding")
		if accept == "gzip" {
			if gzr, err := gzip.NewReader(r.Body); err != nil {
				log.Error("gzip reader error:", err)
			} else {
				r.Body = gzr
			}
		}

		return
	})
}
