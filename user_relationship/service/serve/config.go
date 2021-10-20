package serve

import (
	"github.com/haiyiyun/plugins/user_relationship/service/base"
)

type UserConfig struct {
	WebRouter         bool   `json:"web_router"`
	WebRouterRootPath string `json:"web_router_root_path"`
}

type Config struct {
	base.Config
	UserConfig
}
