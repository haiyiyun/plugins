package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Profile struct {
	UserID     primitive.ObjectID `json:"user_id" bson:"_id" map:"_id"`
	Enable     bool               `json:"enable" bson:"enable" map:"enable"`
	CreateTime time.Time          `json:"create_time" bson:"create_time" map:"create_time"`
	UpdateTime time.Time          `json:"update_time" bson:"update_time" map:"update_time"`
}
