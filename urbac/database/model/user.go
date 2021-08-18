package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSetting struct {
	HomePath string            `bson:"home_path" json:"home_path,omitempty" map:"home_path"`
	Profile  map[string]string `bson:"profile" json:"profile,omitempty" map:"profile"`
	Style    map[string]string `bson:"style" json:"style,omitempty" map:"style"`
}

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"user_id" map:"_id"`
	Name        string             `bson:"name" json:"username" map:"name"`
	RealName    string             `bson:"real_name" json:"real_name" map:"real_name"`
	Email       string             `bson:"email" json:"email" map:"email"`
	Password    string             `bson:"password" json:"-" map:"password"`
	Avatar      string             `bson:"avatar" json:"avatar" map:"avatar"`
	Description string             `bson:"description" json:"description" map:"description"`
	Enable      bool               `bson:"enable" json:"enable" map:"enable"`
	Delete      bool               `bson:"delete" json:"delete" map:"delete"`
	Setting     UserSetting        `bson:"setting" json:"setting" map:"setting"`
	CreateTime  time.Time          `bson:"create_time" json:"create_time" map:"create_time"` //创建时间
	UpdateTime  time.Time          `bson:"update_time" json:"update_time" map:"update_time"` //更新时间
}
