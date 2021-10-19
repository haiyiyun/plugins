package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
"follow_relationship": {
	"_id": "id",
	"create_time": "create_time",
	"update_time": "update_time"
}
*/

var (
	FollowRelationship = bson.M{
		"name":  "follow_relationship",
		"index": []mongo.IndexModel{},
	}
)
