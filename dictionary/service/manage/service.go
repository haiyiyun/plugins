package manage

import (
	"github.com/haiyiyun/plugins/dictionary/service/base"
)

type Service struct {
	*Config
	*base.Service
}

func NewService(c *Config, s *base.Service) *Service {
	return &Service{
		Config:  c,
		Service: s,
	}
}
