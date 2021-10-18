package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
"contacts_blacklist": {
	"_id": "<userid>",
	"create_time": "<create_time>",
	"update_time": "<update_time>"
}
*/

var (
	ContactsBlacklist = bson.M{
		"name": "contacts_blacklist",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"user_id", 1},
					{"blacklist_user_id", 1},
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
					{"blacklist_user_id", 1},
				},
				options.Index(),
			},
		},
	}
)
