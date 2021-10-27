package user

import (
	"context"
	"encoding/gob"

	"github.com/haiyiyun/plugins/urbac/database/model"
	"github.com/haiyiyun/utils/help"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	gob.Register(model.User{})
}

func (self *Model) GetUserByID(userID primitive.ObjectID) (u model.User, err error) {
	filter := self.FilterByID(userID)
	filter = append(filter, bson.D{
		{"enable", true},
		{"delete", false},
	}...)
	sr := self.FindOne(context.TODO(), filter)
	err = sr.Decode(&u)

	return
}

func (self *Model) ChangePassword(userID primitive.ObjectID, password string) error {
	_, err := self.Set(context.TODO(), self.FilterByID(userID), bson.D{
		{"password", help.NewString(password).Md5()},
	})

	return err
}
