package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		"name": "discuss",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"type", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"object_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"publish_user_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"reply_discuss_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"visibility", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"evaluation", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"location", "2dsphere"},
				},
				options.Index(),
			},
			{
				bson.D{
					{"status", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"create_time", -1},
				},
				options.Index(),
			},
		},
	}
)
