package auth

import "github.com/haiyiyun/plugins/urbac/service"

type Service struct {
	*service.Service
}

func NewService(s *service.Service) *Service {
	return &Service{
		Service: s,
	}
}
