package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KeywordBan struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id" map:"_id,omitempty"` //ID
	Type       int
	Keyword    string    `bson:"keyword" json:"keyword" map:"keyword"`             //关键词
	Replace    string    `bson:"replace" json:"replace" map:"replace"`             //替换
	Action     int       `bson:"action" json:"action" map:"action"`                //动作：替换，屏蔽，删除，审核等
	CreateTime time.Time `bson:"create_time" json:"create_time" map:"create_time"` //创建时间
	UpdateTime time.Time `bson:"update_time" json:"update_time" map:"update_time"` //更新时间
}
