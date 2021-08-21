package role

import (
	"github.com/haiyiyun/plugins/urbac/service/manage"
)

type Service struct {
	*manage.Service
}

func NewService(s *manage.Service) *Service {
	return &Service{
		Service: s,
	}
}
