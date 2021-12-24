package serve

import (
	"github.com/haiyiyun/plugins/log/service/base"
)

type LogConfig struct {
	Log               bool   `json:"log"`
	WebRouter         bool   `json:"web_router"`
	WebRouterRootPath string `json:"web_router_root_path"`
}

type Config struct {
	base.Config
	LogConfig
}
