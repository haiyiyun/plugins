package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
"favorites": {
	"_id": "id",
	"create_time": "create_time",
	"update_time": "update_time"
}
*/

var (
	Favorites = bson.M{
		"name":  "favorites",
		"index": []mongo.IndexModel{},
	}
)
