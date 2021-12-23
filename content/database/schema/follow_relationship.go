package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		"name": "follow_relationship",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"type", 1},
					{"user_id", 1},
					{"object_id", 1},
				},
				options.Index().SetUnique(true),
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
					{"user_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"type", 1},
					{"object_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"mutual", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"stealth", 1},
				},
				options.Index(),
			},
		},
	}
)
