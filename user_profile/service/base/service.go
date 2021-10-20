package base

import (
	"github.com/haiyiyun/cache"
	"github.com/haiyiyun/mongodb"
)

type Service struct {
	*Config
	*cache.Cache
	M mongodb.Mongoer
}

func NewService(c *Config, cc *cache.Cache, m mongodb.Mongoer) *Service {
	return &Service{
		Config: c,
		Cache:  cc,
		M:      m,
	}
}
