package base

import (
	"github.com/haiyiyun/config"
)

type MongodbCfg struct {
	MongoDNS          string `json:"mongo_dns"`
	MongoDatabaseName string `json:"mongo_database_name"`
}

type CacheCfg struct {
	CacheType              string          `json:"cache_type"`
	CacheUrl               string          `json:"cache_url"`
	CacheDefaultExpiration config.Duration `json:"cache_default_expiration"`
	CacheCleanupInterval   config.Duration `json:"cache_cleanup_interval"`
}

type Config struct {
	MongodbCfg
	CacheCfg
}
