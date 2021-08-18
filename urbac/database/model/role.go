package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoleRight struct {
	Scope        int                    `bson:"scope" json:"scope" map:"scope"`
	Applications map[string]Application `bson:"applications" json:"applications" map:"applications"`
}

type Role struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty" json:"_id" map:"_id"`
	Name       string               `bson:"name" json:"name" map:"name"`
	Users      []primitive.ObjectID `bson:"users" json:"users" map:"users"`
	Right      RoleRight            `bson:"right" json:"right" map:"right"`
	Enable     bool                 `bson:"enable" json:"enable" map:"enable"`
	Delete     bool                 `bson:"delete" json:"delete" map:"delete"`
	CreateTime time.Time            `bson:"create_time" json:"create_time" map:"create_time"` //创建时间
	UpdateTime time.Time            `bson:"update_time" json:"update_time" map:"update_time"` //更新时间
}
