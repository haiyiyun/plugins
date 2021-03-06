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
    "real_name": "<real_name>",
    "email": "<email>",
    "password": "<password>",
    "avatar": "<avatar>",
    "description": "<description>",
    "status": NumberInt(<status>),
    "delete": NumberInt(<delete>),
    "setting": {
        "profile": {
            "<key>": "<value>"
        },
		"style": {
            "<key>": "<value>"
        }
    },
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
					{"email", 1},
				},
				options.Index().SetUnique(true),
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
