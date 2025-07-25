package base

import (
	"github.com/haiyiyun/config"
)

type MongodbConfig struct {
	MongoDNS          string `json:"mongo_dns"`
	MongoDatabaseName string `json:"mongo_database_name"`
}

type CacheConfig struct {
	CacheType              string          `json:"cache_type"`
	CacheUrl               string          `json:"cache_url"`
	CacheShardCount        string          `json:"cache_shard_count"`
	CacheUStrictTypeCheck  string          `json:"cache_strict_type_check"`
	CacheDefaultExpiration config.Duration `json:"cache_default_expiration"`
	CacheCleanupInterval   config.Duration `json:"cache_cleanup_interval"`
}

type BaseConfig struct {
	URBAC                bool                `json:"urbac"`
	CheckRight           bool                `json:"check_right"`
	TokenByUrlQuery      bool                `json:"token_by_url_query"`
	TokenByUrlQueryName  string              `json:"token_by_url_query_name"`
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

type FrontConfig struct {
	DefaultRoute    map[string]interface{} `json:"default_route"`
	DefaultHomePath string                 `json:"default_home_path"`
}

type Config struct {
	MongodbConfig
	CacheConfig
	BaseConfig
	FrontConfig
}
