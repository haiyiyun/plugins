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
	CheckLogin           bool                `json:"check_login"`
	TokenByUrlQuery      bool                `json:"token_by_url_query"`
	TokenByUrlQueryName  string              `json:"token_by_url_query_name"`
	IgnoreCheckLoginPath map[string][]string `json:"ignore_check_login_path"`
	TokenExpireDuration  config.Duration     `json:"token_expire_duration"`
	AllowMultiLogin      bool                `json:"allow_multi_login"`
	AllowMultiLoginNum   int64               `json:"allow_multi_login_num"`
}

type Config struct {
	MongodbCfg
	CacheCfg
	BaseCfg
}
