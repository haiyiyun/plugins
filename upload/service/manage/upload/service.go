package upload

import "github.com/haiyiyun/plugins/upload/service/manage"

type Service struct {
	*manage.Service
}

func NewService(s *manage.Service) *Service {
	return &Service{
		Service: s,
	}
}
