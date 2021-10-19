package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
"publish_intervene": {
	"_id": "id",
	"create_time": "create_time",
	"update_time": "update_time"
}
*/

var (
	PublishIntervene = bson.M{
		"name": "publish_intervene",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"user_id", 1},
					{"type", 1},
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
