package dictionary

import "github.com/haiyiyun/plugins/dictionary/service/serve"

type Service struct {
	*serve.Service
}

func NewService(s *serve.Service) *Service {
	return &Service{
		Service: s,
	}
}
