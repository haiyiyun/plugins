package manage

import (
	"github.com/haiyiyun/plugins/dictionary/service/base"
)

type DictionaryCfg struct {
	WebRouter         bool   `json:"web_router"`
	WebRouterRootPath string `json:"web_router_root_path"`
}

type Config struct {
	base.Config
	DictionaryCfg
}
