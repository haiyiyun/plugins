package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Contacts struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id" map:"_id"`
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id" map:"user_id"`
	ContactUserID primitive.ObjectID `bson:"contact_user_id" json:"contact_user_id" map:"contact_user_id"`
	CreateTime    time.Time          `bson:"create_time" json:"create_time" map:"create_time"` //创建时间
	UpdateTime    time.Time          `bson:"update_time" json:"update_time" map:"update_time"` //更新时间
}
