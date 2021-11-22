package user

import (
	"context"
	"encoding/gob"

	"github.com/haiyiyun/plugins/user/database/model"
	"github.com/haiyiyun/utils/help"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	gob.Register(model.User{})
}

func (self *Model) FilterByName(name string) bson.D {
	return bson.D{
		{"name", name},
	}
}

func (self *Model) CheckNameAndPassword(name, password string) (u model.User, err error) {
	passwordMd5 := help.Strings(password).Md5()

	filter := bson.D{
		{"name", name},
		{"password", passwordMd5},
		{"enable", true},
		{"delete", false},
	}

	sr := self.FindOne(context.TODO(), filter)
	err = sr.Decode(&u)

	return
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

func (self *Model) GuestToUser(userID primitive.ObjectID, username, password string) error {
	_, err := self.Set(context.TODO(), self.FilterByID(userID), bson.D{
		{"name", username},
		{"password", help.NewString(password).Md5()},
		{"guest", false},
	})

	return err
}
