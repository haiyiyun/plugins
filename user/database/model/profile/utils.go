package profile

import (
	"context"

	"github.com/haiyiyun/utils/help"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (self *Model) GetInfo(userID primitive.ObjectID) (pf help.M, err error) {
	filter := self.FilterByID(userID)
	filter = append(filter, bson.D{
		{"enable", true},
	}...)

	sr := self.FindOne(context.TODO(), filter)
	err = sr.Decode(&pf)

	return
}
