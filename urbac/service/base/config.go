package base

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

type BaseCfg struct {
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
}

type FrontCfg struct {
	DefaultRoute    map[string]interface{} `json:"default_route"`
	DefaultHomePath string                 `json:"default_home_path"`
}

type Config struct {
	MongodbCfg
	CacheCfg
	BaseCfg
	FrontCfg
}
