package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
"contacts": {
	"_id": "<userid>",
	"status": "<status>",
	"create_time": "<create_time>",
	"update_time": "<update_time>"
}
*/

var (
	Contacts = bson.M{
		"name": "contacts",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"user_id", 1},
					{"contact_user_id", 1},
				},
				options.Index().SetUnique(true),
			},
			{
				bson.D{
					{"user_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"contact_user_id", 1},
				},
				options.Index(),
			},
		},
	}
)
