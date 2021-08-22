package manage

import (
	"github.com/haiyiyun/plugins/urbac/service/base"
)

type UrbacConfig struct {
	WebRouter         bool   `json:"web_router"`
	WebRouterRootPath string `json:"web_router_root_path"`
}

type Config struct {
	base.Config
	UrbacConfig
}
