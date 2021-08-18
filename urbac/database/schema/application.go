package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
应用表
application
{
    "_id": "<_id>",
    "name": "<name>",
    "icon": "<icon>",
    "path": "<path>", # /baisc/
    "order": NumberLong(<order>),
    "modules": {
        "<path>": { # user/
            "name": "<name>", # 用户管理
            "icon": "<icon>", # &#xe723;
            "path": "<path>", # user/
            "order": NumberLong(<order>),
            "actions": {
                "<path>": { # index_menu
                    "name": "<name>", # 用户列表
                    "icon": "<icon>", # &#xe6a7;
                    "path": "<path>", # index_menu
                    "order": NumberLong(<order>)
                }
            }
        }
    },
    "status": NumberInt(<status>),
    "create_time": "<create_time>", # ISODate("2019-08-15T07:31:32.845Z"),
    "update_time": "<update_time>" # ISODate("2019-08-15T07:31:32.845Z")
}
*/

var (
	Application = bson.M{
		"name": "application",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"name", 1},
				},
				options.Index().SetUnique(true),
			},
			{
				bson.D{
					{"path", 1},
				},
				options.Index().SetUnique(true),
			},
			{
				bson.D{
					{"enable", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"order", 1},
				},
				options.Index(),
			},
		},
	}
)
