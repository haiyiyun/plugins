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
	DefaultDeleteDuration        config.Duration `json:"default_delete_duration"`
	DefaultLoginDeleteDuration   config.Duration `json:"default_login_delete_duration"`
	DefaultAuthDeleteDuration    config.Duration `json:"default_auth_delete_duration"`
	DefaultOperateDeleteDuration config.Duration `json:"default_operate_delete_duration"`
	LogLoginPath                 []string        `json:"log_login_path"`
	LogAuthPath                  []string        `json:"log_auth_path"`
	LogOperateExcludePath        []string        `json:"log_operate_exclude_path"`
}

type Config struct {
	MongodbCfg
	CacheCfg
	BaseCfg
}
