package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
"group_message": {
	"_id": "id",
	"create_time": "create_time",
	"update_time": "update_time"
}
*/

var (
	GroupMessage = bson.M{
		"name":  "group_message",
		"index": []mongo.IndexModel{},
	}
)
