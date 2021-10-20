package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Profile = bson.M{
		"name":  "profile",
		"index": []mongo.IndexModel{},
	}
)
