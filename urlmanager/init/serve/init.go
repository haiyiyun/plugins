package serve

import (
	cmdflag "flag"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/plugins/urlmanager/service/serve"
	"github.com/haiyiyun/webrouter"
)

func init() {
	urlManagerConfFile := cmdflag.String("config.plugins.urlmanager", "../config/plugins/urlmanager/urlmanager.conf", "urlmanager config file")
	config, err := ioutil.ReadFile(filepath.Clean(*urlManagerConfFile))
	if err != nil {
		log.Error("readfile config.plugins.urlmanager faild:", err)
		os.Exit(-1)
	}

	urlManager := serve.New()
	urlManager.Start()
	urlManager.LoadRule(string(config), false)

	webrouter.Injector("urlmanager", "", 9999, func(w http.ResponseWriter, r *http.Request) (abort bool) {
		if urlManager.Manage() {
			newUrl := urlManager.ReWrite(w, r)
			if newUrl == "redirect" {
				abort = true
			} else {
				r.URL, _ = url.Parse(newUrl)
			}
		}

		return
	})
}
