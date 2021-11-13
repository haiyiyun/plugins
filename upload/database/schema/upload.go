package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
上传文件表
upload
{
	"_id": "<ID>",
    "type": "<Type>",
    "storage": "<Storage>",
    "user_id": "<UserID>",
    "file_name": "<FileName>",
    "path": "<Path>",
    "url": "<url>",
    "size": "<size>",
    "create_time": "<create_time>", # ISODate("2019-08-19T01:22:44.234Z")
    "update_time": "<update_time>" # ISODate("2019-08-19T01:22:44.234Z")
}
*/

var (
	Upload = bson.M{
		"name": "upload",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"type", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"storage", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"user_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"content_type", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"original_file_name", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"file_name", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"file_ext", 1},
				},
				options.Index(),
			},
		},
	}
)
