package model

type Country struct {
	ID   string `json:"_id" bson:"_id" map:"_id"`
	Name string `json:"name" bson:"name" map:"name"`
}
