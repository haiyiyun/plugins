package upload

import (
	"github.com/haiyiyun/plugins/upload/service/serve"
)

type Service struct {
	*serve.Service
}

func NewService(s *serve.Service) *Service {
	return &Service{
		Service: s,
	}
}
