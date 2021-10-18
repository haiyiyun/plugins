package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
用户表
user
{
	"_id": "<_id>",
    "name": "<name>",
    "password": "<password>",
    "status": NumberInt(<status>),
    "delete": NumberInt(<delete>),
    "create_time": "<create_time>", # ISODate("2019-08-19T01:22:44.234Z")
    "update_time": "<update_time>" # ISODate("2019-08-19T01:22:44.234Z")
}
*/

var (
	User = bson.M{
		"name": "user",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"name", 1},
				},
				options.Index().SetUnique(true),
			},
			{
				bson.D{
					{"guest", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"level", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"role.role", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"online.online", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"online.location", "2dsphere"},
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
					{"enable", 1},
					{"delete", 1},
				},
				options.Index(),
			},
			//来宾模式下，7天后自动删除
			// {
			// 	bson.D{
			// 		{"create_time", 1},
			// 	},
			// 	options.Index().SetPartialFilterExpression(bson.D{
			// 		{"guest", true},
			// 	}).SetExpireAfterSeconds(60 * 60 * 24 * 7),
			// },
		},
	}
)
