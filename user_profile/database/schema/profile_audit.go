package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
"profile_audit": {
	"_id": "<userid>",
	"status": "<status>",
	"create_time": "<create_time>",
	"update_time": "<update_time>"
}
*/

var (
	ProfileAudit = bson.M{
		"name": "profile_audit",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"type", 1},
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
