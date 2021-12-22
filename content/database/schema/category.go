package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
"category": {
	"_id": "id",
	"create_time": "create_time",
	"update_time": "update_time"
}
*/

var (
	Category = bson.M{
		"name": "category",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"type", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"parent_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"tags", 1},
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
					{"visibility", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"limit_user_at_least_level", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"only_user_id_not_limit_user_level", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"limit_user_role", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"only_user_id_not_limit_user_role", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"limit_user_tag", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"only_user_id_not_limit_user_tag", 1},
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
					{"enable", 1},
				},
				options.Index(),
			},
		},
	}
)
