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
	WebRouter         bool   `json:"web_router"`
	WebRouterRootPath string `json:"web_router_root_path"`
}

type Config struct {
	MongodbConfig
	CacheConfig
	BaseConfig
}
