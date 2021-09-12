package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
"country": {
	"_id": "id",
	"name": "name",
}
*/

var (
	Country = bson.M{
		"name":  "country",
		"index": []mongo.IndexModel{},
	}
)
