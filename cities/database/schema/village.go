package schema

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
"village": {
	"_id": "id",
	"name": "name",
	"province_id": "province_id",
	"city_id": "city_id",
	"area_id": "area_id",
	"street_id": "street_id",
}
*/

var (
	Village = bson.M{
		"name": "village",
		"index": []mongo.IndexModel{
			{
				bson.D{
					{"province_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"city_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"area_id", 1},
				},
				options.Index(),
			},
			{
				bson.D{
					{"street_id", 1},
				},
				options.Index(),
			},
		},
	}
)
