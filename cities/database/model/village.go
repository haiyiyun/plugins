package model

type Village struct {
	ID         string `json:"_id" bson:"_id" map:"_id"`
	Name       string `json:"name" bson:"name" map:"name"`
	ProvinceID string `json:"province_id" bson:"province_id" map:"province_id"`
	CityID     string `json:"city_id" bson:"city_id" map:"city_id"`
	AreaID     string `json:"area_id" bson:"area_id" map:"area_id"`
	StreetID   string `json:"street_id" bson:"street_id" map:"street_id"`
}
