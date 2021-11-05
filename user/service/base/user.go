package base

import (
	"context"
	"net/http"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/mongodb/geometry"
	"github.com/haiyiyun/plugins/user/database/model"
	"github.com/haiyiyun/plugins/user/database/model/profile"
	"github.com/haiyiyun/plugins/user/database/model/user"
	"github.com/haiyiyun/utils/help"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (self *Service) CreateUser(ctx context.Context, username, password string, longitude, latitude float64, enableProfile bool) (primitive.ObjectID, error) {
	coordinates := geometry.PointCoordinates{
		longitude, latitude,
	}

	userModel := user.NewModel(self.M)
	var userID primitive.ObjectID
	err := userModel.UseSession(ctx, func(sctx mongo.SessionContext) error {
		u := model.User{
			Name:     username,
			Password: help.NewString(password).Md5(),
			Enable:   true,
		}

		if coordinates != geometry.NilPointCoordinates {
			u.Location = geometry.NewPoint(coordinates)
		}

		ior, err := userModel.Create(sctx, u)

		if err != nil {
			sctx.AbortTransaction(sctx)
			log.Error("Create user error:", err)
			return err
		}

		userID = ior.InsertedID.(primitive.ObjectID)

		if enableProfile {
			profileModel := profile.NewModel(self.M)
			_, err = profileModel.Create(sctx, model.Profile{
				UserID: userID,
				Enable: true,
			})

			if err != nil {
				sctx.AbortTransaction(sctx)
				log.Error("Create profile error:", err)
				return err
			}
		}

		sctx.CommitTransaction(sctx)
		return err
	})

	return userID, err
}
