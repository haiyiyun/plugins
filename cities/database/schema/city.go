package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
"city": {
	"_id": "id",
	"name": "name",
	"province_id": "province_id",
}
*/

var (
	City = bson.M{
		"name": "city",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"province_id", 1},
				},
				options.Index(),
			},
		},
	}
)
