package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//阅读过后的一段时间删除访客记录，比如：一天后
type Visitor struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id" map:"_id"`
	OwnerUserID   primitive.ObjectID `bson:"owner_user_id" json:"owner_user_id" map:"owner_user_id"`
	VisitorUserID primitive.ObjectID `bson:"visitor_user_id" json:"visitor_user_id" map:"visitor_user_id"`
	Hide          bool               `bson:"hide" json:"hide" map:"hide"`                         //是否隐藏访问
	BeShowTime    time.Time          `bson:"be_show_time" json:"be_show_time" map:"be_show_time"` //被显示时间
	ReadedTime    time.Time          `bson:"readed_time" json:"readed_time" map:"readed_time"`
	CreateTime    time.Time          `bson:"create_time" json:"create_time" map:"create_time"` //创建时间
	UpdateTime    time.Time          `bson:"update_time" json:"update_time" map:"update_time"` //更新时间
}
