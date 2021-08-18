package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
用户表
log
{
	"_id": "<_id>",
    "type": "<type>",
    "user": "<user>",
    "method": "<method>",
    "path": "<path>",
    "query": "<query>",
    "ip": "<ip>",
    "create_time": "<create_time>", # ISODate("2019-08-19T01:22:44.234Z")
    "update_time": "<update_time>" # ISODate("2019-08-19T01:22:44.234Z")
}
*/

var (
	Log = bson.M{
		"name": "log",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"type", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"user", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"method", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"path", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"delete_time", 1},
				},
				options.Index().SetExpireAfterSeconds(1),
			},
		},
	}
)
