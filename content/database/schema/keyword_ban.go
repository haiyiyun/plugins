package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
"keyword_ban": {
	"_id": "id",
	"create_time": "create_time",
	"update_time": "update_time"
}
*/

var (
	KeywordBan = bson.M{
		"name":  "keyword_ban",
		"index": []mongo.IndexModel{},
	}
)
