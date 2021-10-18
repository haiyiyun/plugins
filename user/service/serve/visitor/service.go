package visitor

import "github.com/haiyiyun/plugins/user/service/serve"

type Service struct {
	*serve.Service
}

func NewService(s *serve.Service) *Service {
	return &Service{
		Service: s,
	}
}
