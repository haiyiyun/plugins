package manage

import (
	"github.com/haiyiyun/plugins/user_profile/service/base"
)

type WebrouterPluginTemplateConfig struct {
	WebRouter         bool   `json:"web_router"`
	WebRouterRootPath string `json:"web_router_root_path"`
}

type Config struct {
	base.Config
	WebrouterPluginTemplateConfig
}
