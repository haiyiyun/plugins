package log

import "github.com/haiyiyun/plugins/log/service/manage"

type Service struct {
	*manage.Service
}

func NewService(s *manage.Service) *Service {
	return &Service{
		Service: s,
	}
}
