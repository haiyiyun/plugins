package role

import (
	"context"
	"encoding/gob"

	"github.com/haiyiyun/plugins/urbac/database/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	gob.Register(model.Role{})
	gob.Register([]model.Role{})
}

func (self *Model) FilterByName(name string) bson.D {
	return bson.D{
		{"name", name},
	}
}

func (self *Model) GetRoleByUserID(userID primitive.ObjectID) (roles []model.Role, err error) {
	filter := bson.D{
		{"users", userID},
		{"enable", true},
		{"delete", false},
	}

	opt := options.Find().SetSort(bson.D{
		{"right.scope", -1},
	})

	cur := new(mongo.Cursor)
	if cur, err = self.Find(context.Background(), filter, opt); err == nil {
		err = cur.All(context.TODO(), &roles)
	}

	return
}
