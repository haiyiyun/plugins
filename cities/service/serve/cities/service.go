package cities

import (
	"github.com/haiyiyun/plugins/cities/service/serve"
)

type Service struct {
	*serve.Service
}

func NewService(s *serve.Service) *Service {
	return &Service{
		Service: s,
	}
}
