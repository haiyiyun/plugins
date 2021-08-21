package manage

import (
	"github.com/haiyiyun/plugins/urbac/service/base"
)

type UrbacCfg struct {
	WebRouter         bool   `json:"web_router"`
	WebRouterRootPath string `json:"web_router_root_path"`
}

type Config struct {
	base.Config
	UrbacCfg
}
