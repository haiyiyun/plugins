package dictionary

import "github.com/haiyiyun/plugins/dictionary/service/manage"

type Service struct {
	*manage.Service
}

func NewService(s *manage.Service) *Service {
	return &Service{
		Service: s,
	}
}
