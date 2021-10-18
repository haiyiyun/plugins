package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
"visitor": {
	"_id": "<userid>",
	"create_time": "<create_time>",
	"update_time": "<update_time>"
}
*/

var (
	Visitor = bson.M{
		"name": "visitor",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"owner_user_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"visitor_user_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"hide", 1},
				},
				options.Index(),
			},
		},
	}
)
