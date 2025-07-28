package upload

import (
	"github.com/haiyiyun/log"
	"github.com/haiyiyun/plugins/upload/service"
	"github.com/haiyiyun/plugins/upload/service/serve"
)

type Service struct {
	*serve.Service
	uploadService service.Upload // 使用统一的Upload接口
}

func NewService(s *serve.Service) *Service {
	// 通过工厂创建存储服务
	uploadService, err := service.NewUploadService(s.Service)
	if err != nil {
		log.Fatal("Failed to create upload service:", err)
	}

	return &Service{
		Service:       s,
		uploadService: uploadService,
	}
}
