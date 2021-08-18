package service

import (
	"net/http"

	"github.com/haiyiyun/plugins/urbac/database/model"
	"github.com/haiyiyun/plugins/urbac/database/model/role"
	"github.com/haiyiyun/plugins/urbac/database/model/user"

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

func (self *Service) getRole(userID primitive.ObjectID) (roles []model.Role, err error) {
	roleModel := role.NewModel(self.M)
	roles, err = roleModel.GetRoleByUserID(userID)

	return
}
