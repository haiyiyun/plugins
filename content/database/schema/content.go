package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
"content": {
	"_id": "id",
	"create_time": "create_time",
	"update_time": "update_time"
}
*/

var (
	Content = bson.M{
		"name": "content",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"publish_user_id", 1},
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
					{"publish_type", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"associate_type", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"associate_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"category_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"subject_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"at_users", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"author", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"location", "2dsphere"},
				},
				options.Index(),
			},
			{
				bson.D{
					{"publish_user_id", 1},
					{"user_tags", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"visibility", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"copy", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"forbid_forward", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"forbid_download", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"forbid_discuss", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"tags", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"bestest", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"reliable", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"discuss_estimate_total", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"value", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"start_time", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"end_time", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"status", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"create_time", -1},
				},
				options.Index(),
			},
		},
	}
)
