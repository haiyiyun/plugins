package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		"name": "follow_content",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"follow_relationship_id", 1},
					{"user_id", 1},
					{"content_id", 1},
				},
				options.Index().SetUnique(true),
			},
			{
				bson.D{
					{"user_id", 1},
					{"content_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"follow_relationship_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"user_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"content_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"extension_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"type", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"readed_time", 1},
				},
				options.Index(),
			},
			// {
			// 	bson.D{
			// 		{"readed_time", 1},
			// 	},
			// 	options.Index().SetExpireAfterSeconds(60 * 60 * 24 * 30),
			// },
		},
	}
)
