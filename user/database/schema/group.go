package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
"group": {
	"_id": "<userid>",
	"status": "<status>",
	"create_time": "<create_time>",
	"update_time": "<update_time>"
}
*/

var (
	Group = bson.M{
		"name": "group",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"_id", 1},
					{"members.user_id", 1},
				},
				options.Index().SetUnique(true),
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
					{"hide", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"hide_members", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"join", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"audit", 1},
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
					{"delete", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"members.adminer", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"members.stealth", 1},
				},
				options.Index(),
			},
		},
	}
)
