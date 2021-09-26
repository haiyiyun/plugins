package pre_parseform

import (
	"net/http"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/utils/http/request"
	"github.com/haiyiyun/webrouter"
)

func init() {
	webrouter.Injector("pre_parseform", "", 99, func(rw http.ResponseWriter, req *http.Request) (abort bool) {
		switch req.Method {
		case "GET":
			switch req.Header.Get("Content-Type") {
			case "application/x-www-form-urlencoded":
				if bBody, err := request.GetBody(req); err != nil {
					log.Error("ReadAll body error:", err)
				} else {
					req.URL.RawQuery = string(bBody)
				}
			}
		case "POST", "PUT":
			switch req.Header.Get("Content-Type") {
			case "application/x-www-form-urlencoded":
				req.ParseForm()
			}
		}

		return
	})
}
