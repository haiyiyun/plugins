package serve

import (
	"github.com/haiyiyun/plugins/dictionary/service/base"
)

type DictionaryConfig struct {
	WebRouter         bool   `json:"web_router"`
	WebRouterRootPath string `json:"web_router_root_path"`
}

type Config struct {
	base.Config
	DictionaryConfig
}
