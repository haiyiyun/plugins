package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
"discuss": {
	"_id": "id",
	"create_time": "create_time",
	"update_time": "update_time"
}
*/

var (
	Discuss = bson.M{
		"name":  "discuss",
		"index": []mongo.IndexModel{},
	}
)
