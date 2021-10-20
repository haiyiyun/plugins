package serve

import "github.com/haiyiyun/plugins/user_relationship/service/base"

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
