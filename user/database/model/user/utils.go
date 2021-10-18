package user

import (
	"context"
	"encoding/gob"

	"github.com/haiyiyun/plugins/user/database/model"

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
