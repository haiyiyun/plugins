package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
"follow_content": {
	"_id": "id",
	"create_time": "create_time",
	"update_time": "update_time"
}
*/

var (
	FollowContent = bson.M{
		"name":  "follow_content",
		"index": []mongo.IndexModel{},
	}
)
