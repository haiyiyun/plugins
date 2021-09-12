package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
"province": {
	"_id": "id",
	"name": "name",
}
*/

var (
	Province = bson.M{
		"name":  "province",
		"index": []mongo.IndexModel{},
	}
)
