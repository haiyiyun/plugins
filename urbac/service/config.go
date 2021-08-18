package service

import (
	"github.com/haiyiyun/config"
)

type MongodbCfg struct {
	MongoDNS          string `json:"mongo_dns"`
	MongoDatabaseName string `json:"mongo_database_name"`
}

type CacheCfg struct {
	CacheDefaultExpiration config.Duration `json:"cache_default_expiration"`
	CacheCleanupInterval   config.Duration `json:"cache_cleanup_interval"`
}

type UrbacCfg struct {
	WebRouter            bool                `json:"web_router"`
	WebRouterRootPath    string              `json:"web_router_root_path"`
	CheckRight           bool                `json:"check_right"`
	DefaultEnableApp     bool                `json:"default_enable_app"`
	DefaultEnableModule  bool                `json:"default_enable_module"`
	DefaultEnableAction  bool                `json:"default_enable_action"`
	IgnoreAppModuleInfo  []string            `json:"ignore_app_module_info"`
	IgnoreCheckLoginPath map[string][]string `json:"ignore_check_login_path"`
	IgnoreCheckRightPath map[string][]string `json:"ignore_check_right_path"`
	TokenExpireDuration  config.Duration     `json:"token_expire_duration"`
	AllowMultiLogin      bool                `json:"allow_multi_login"`
	AllowMultiLoginNum   int64               `json:"allow_multi_login_num"`
	DefaultHomePath      string              `json:"default_home_path"`
}

type FrontCfg struct {
	DefaultRoute map[string]interface{} `json:"default_route"`
}

type Config struct {
	MongodbCfg
	CacheCfg
	UrbacCfg
	FrontCfg
}
