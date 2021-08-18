package upload

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (self *Model) FilterByUserID(userID primitive.ObjectID) bson.D {
	return bson.D{
		{"user_id", userID},
	}
}
