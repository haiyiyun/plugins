package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
"audit": {
	"_id": "<userid>",
	"enable": "<enable>",
	"create_time": "<create_time>",
	"update_time": "<update_time>"
}
*/

var (
	Audit = bson.M{
		"name": "audit",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"type", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"enable", 1},
				},
				options.Index(),
			},
		},
	}
)
