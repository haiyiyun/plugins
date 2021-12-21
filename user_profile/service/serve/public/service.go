package public

import "github.com/haiyiyun/plugins/user_profile/service/serve"

type Service struct {
	*serve.Service
}

func NewService(s *serve.Service) *Service {
	return &Service{
		Service: s,
	}
}
