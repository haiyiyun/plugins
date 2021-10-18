package base

import (
	"net/http"

	"github.com/haiyiyun/plugins/user/database/model"
	"github.com/haiyiyun/plugins/user/database/model/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (self *Service) GetUserInfo(r *http.Request) (u model.User, found bool) {
	if claims, uTmp := self.GetValidClaims(r); claims != nil {
		u = uTmp
		found = true
	}

	return
}

func (self *Service) getUser(userID primitive.ObjectID) (u model.User, err error) {
	userModel := user.NewModel(self.M)
	u, err = userModel.GetUserByID(userID)

	return
}
