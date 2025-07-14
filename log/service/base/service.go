package base

import (
	"encoding/gob"
	"time"

	"github.com/haiyiyun/cache"
	"github.com/haiyiyun/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	gob.Register(map[string]bool{})
	gob.Register(map[string]int{})
	gob.Register(map[string]string{})
	gob.Register(map[string]interface{}{})
	gob.Register([]map[string]interface{}{})
	gob.Register(time.Time{})
	gob.Register([]primitive.ObjectID{})
}

type Service struct {
	*Config
	cache.Cache
	M mongodb.Mongoer
}

func NewService(c *Config, cc cache.Cache, m mongodb.Mongoer) *Service {
	return &Service{
		Config: c,
		Cache:  cc,
		M:      m,
	}
}
