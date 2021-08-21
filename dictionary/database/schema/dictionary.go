package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
dictionary
{
	"_id": "<_id>",
    "create_time": "<create_time>", # ISODate("2019-08-19T01:22:44.234Z")
    "update_time": "<update_time>" # ISODate("2019-08-19T01:22:44.234Z")
}
*/

var (
	Dictionary = bson.M{
		"name": "dictionary",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"key", 1},
				},
				options.Index().SetUnique(true),
			},
			{
				bson.D{
					{"key", 1},
					{"values.key", 1},
				},
				options.Index().SetUnique(true),
			},
			{
				bson.D{
					{"key", 1},
					{"values.lable", 1},
				},
				options.Index().SetUnique(true),
			},
			{
				bson.D{
					{"key", 1},
					{"values.value", 1},
				},
				options.Index().SetUnique(true),
			},
			{
				bson.D{
					{"values.enable", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"enable", 1},
					{"delete", 1},
				},
				options.Index(),
			},
		},
	}
)
