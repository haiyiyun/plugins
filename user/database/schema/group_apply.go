package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
"group_apply": {
	"_id": "<userid>",
	"create_time": "<create_time>",
	"update_time": "<update_time>"
}
*/

var (
	GroupApply = bson.M{
		"name": "group_apply",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"group_id", 1},
					{"applyer_user_id", 1},
				},
				options.Index().SetUnique(true),
			},
			{
				bson.D{
					{"group_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"applyer_user_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"refuse", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"pass", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"delete", 1},
				},
				options.Index(),
			},
		},
	}
)
