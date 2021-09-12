package base

import (
	"github.com/haiyiyun/config"
)

type MongodbConfig struct {
	MongoDNS          string `json:"mongo_dns"`
	MongoDatabaseName string `json:"mongo_database_name"`
}

type CacheConfig struct {
	CacheDefaultExpiration config.Duration `json:"cache_default_expiration"`
	CacheCleanupInterval   config.Duration `json:"cache_cleanup_interval"`
}

type BaseConfig struct {
}

type Config struct {
	MongodbConfig
	CacheConfig
	BaseConfig
}
