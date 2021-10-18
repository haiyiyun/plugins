package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
token表
token
{
	"_id": "<_id>",
	"userid": "<userid>",
    "type": "<type>",
    "token": "<token>",
    "expired_time": "<expired_time>" # ISODate("2019-08-19T01:22:44.234Z")
    "create_time": "<create_time>", # ISODate("2019-08-19T01:22:44.234Z")
    "update_time": "<update_time>" # ISODate("2019-08-19T01:22:44.234Z")
}
*/

var (
	Token = bson.M{
		"name": "token",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"userid", 1},
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
					{"sign_info.location", "2dsphere"},
				},
				options.Index(),
			},
			{
				bson.D{
					{"expired_time", 1},
				},
				options.Index().SetExpireAfterSeconds(1),
			},
		},
	}
)
