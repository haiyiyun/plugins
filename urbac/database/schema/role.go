package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
角色表
{
    "_id": "<_id>",
    "name": "<name>",
    "users": [
		"<user_name>"
	],
    "status": NumberInt(<status>),
    "right": {
        "scope": "<scope>", # 0=>action,1=>module,2=>app,3=>site
        "applications": {
            "<path>": {
                "path": "<path>",
                "modules": {
                    "<path>": {
                        "path": "<path>",
                        "actions": {
                            "<path>": {
                                "path": "<path>",
                            }
                        }
                    }
                },
                "status": NumberInt(<status>),
                "create_time": <create_time>, # ISODate("01-01-01T00:00:00.000Z"),
                "update_time": <update_time>, # ISODate("01-01-01T00:00:00.000Z")
            }
        }
    },
    "delete": NumberInt(<delete>),
    "create_time": <create_time>, # ISODate("2019-08-20T15:12:13.844Z"),
    "update_time": <update_time>, # ISODate("2019-08-20T15:35:40.433Z")
}
*/

var (
	Role = bson.M{
		"name": "role",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"name", 1},
				},
				options.Index().SetUnique(true),
			},
			{
				bson.D{
					{"users", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"right.scope", 1},
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
