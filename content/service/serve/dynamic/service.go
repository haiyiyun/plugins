package dynamic

import (
	"github.com/haiyiyun/plugins/content/service/serve"
)

type Service struct {
	*serve.Service
}

func NewService(s *serve.Service) *Service {
	return &Service{
		Service: s,
	}
}
