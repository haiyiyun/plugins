package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupBlacklist struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"_id" map:"_id"`
	GroupID         primitive.ObjectID `bson:"group_id" json:"group_id" map:"group_id"`
	BlacklistUserID primitive.ObjectID `bson:"blacklist_user_id" json:"blacklist_user_id" map:"blacklist_user_id"`
	CreateTime      time.Time          `bson:"create_time" json:"create_time" map:"create_time"` //创建时间
	UpdateTime      time.Time          `bson:"update_time" json:"update_time" map:"update_time"` //更新时间
}
