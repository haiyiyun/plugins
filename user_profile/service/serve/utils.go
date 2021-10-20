package serve

import (
	"net/http"

	"github.com/haiyiyun/utils/http/request"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (self *Service) GetUserID(r *http.Request) primitive.ObjectID {
	claims := request.GetClaims(r)
	if claims == nil {
		return primitive.NilObjectID
	}

	userID, _ := primitive.ObjectIDFromHex(claims.Audience)

	return userID
}
