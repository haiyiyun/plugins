package serve

import (
	"github.com/haiyiyun/plugins/user/service/base"
)

type UserConfig struct {
	WebRouter         bool   `json:"web_router"`
	WebRouterRootPath string `json:"web_router_root_path"`
	EnableProfile     bool   `json:"enable_profile"`
}

type Config struct {
	base.Config
	UserConfig
}
