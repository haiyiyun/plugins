package model

type City struct {
	ID         string `json:"_id" bson:"_id" map:"_id"`
	Name       string `json:"name" bson:"name" map:"name"`
	ProvinceID string `json:"province_id" bson:"province_id" map:"province_id"`
}
