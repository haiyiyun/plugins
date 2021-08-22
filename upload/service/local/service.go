package local

import (
	"github.com/haiyiyun/plugins/upload/service/base"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	*base.Service
	userID primitive.ObjectID
}

func NewService(s *base.Service) *Service {
	return &Service{
		Service: s,
		userID:  primitive.NilObjectID,
	}
}
