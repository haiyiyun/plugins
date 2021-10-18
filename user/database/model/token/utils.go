package token

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (self *Model) FilterByUserID(userID primitive.ObjectID) bson.D {
	return bson.D{
		{"user_id", userID},
	}
}

func (self *Model) FilterByType(typ string) bson.D {
	return bson.D{
		{"token_type", typ},
	}
}

func (self *Model) CountDocumentsByID(jwtID primitive.ObjectID) (int64, error) {
	return self.CountDocuments(context.TODO(), self.FilterByID(jwtID))
}

func (self *Model) CountDocumentsByIDAndToken(jwtID primitive.ObjectID, token string) (int64, error) {
	filter := self.FilterByID(jwtID)
	filter = append(filter, bson.D{
		{"token", token},
	}...)
	return self.CountDocuments(context.TODO(), filter)
}

func (self *Model) CountDocumentsByUserIDAndType(userID primitive.ObjectID, typ string) (int64, error) {
	filter := self.FilterByUserID(userID)
	if typ != "" {
		filter = append(filter, self.FilterByType(typ)...)
	}
	return self.CountDocuments(context.TODO(), filter)
}
