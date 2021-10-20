package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
"profile": {
	"_id": "<userid>",
	"enable": "<enable>",
	"create_time": "<create_time>",
	"update_time": "<update_time>"
}
*/

var (
	Profile = bson.M{
		"name": "profile",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"info.nickname", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"info.sex", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"info.address.type", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"info.address.nation", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"info.address.province", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"info.address.city", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"info.address.district", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"status", 1},
				},
				options.Index(),
			},
		},
	}
)
