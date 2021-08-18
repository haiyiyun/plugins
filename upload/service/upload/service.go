package upload

import "github.com/haiyiyun/plugins/upload/service"

type Service struct {
	*service.Service
}

func NewService(s *service.Service) *Service {
	return &Service{
		Service: s,
	}
}
