package base

import (
	"context"
	"net/http"

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

func (self *Service) CheckUser(ctx context.Context, username string) (int64, error) {
	userModel := user.NewModel(self.M)
	return userModel.CountDocuments(ctx, userModel.FilterByName(username))
}

func (self *Service) CreateUser(ctx context.Context, id primitive.ObjectID, username, password string, longitude, latitude float64, enableProfile bool, extensionID int, guest bool) (primitive.ObjectID, error) {
	coordinates := geometry.PointCoordinates{
		longitude, latitude,
	}

	userModel := user.NewModel(self.M)

	var userID primitive.ObjectID
	err := userModel.UseSession(ctx, func(sctx mongo.SessionContext) error {
		if err := sctx.StartTransaction(); err != nil {
			return err
		}

		u := model.User{
			ExtensionID: extensionID,
			Name:        username,
			Password:    help.NewString(password).Md5(),
			Roles:       []model.UserRole{},
			Tags:        []model.UserTag{},
			Guest:       guest,
			Enable:      true,
		}

		if id != primitive.NilObjectID {
			u.ID = id
		}

		if coordinates != geometry.NilPointCoordinates {
			u.Location = geometry.NewPoint(coordinates)
		}

		ior, err := userModel.Create(sctx, u)

		if err != nil {
			sctx.AbortTransaction(sctx)
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
				return err
			}
		}

		return sctx.CommitTransaction(sctx)
	})

	return userID, err
}
