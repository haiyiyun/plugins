package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProfileAudit struct {
	ID         primitive.ObjectID  `bson:"_id,omitempty" json:"_id" map:"_id"`
	Type       int                 `bson:"type" json:"type" map:"type"`                      //类型，如profile的头像等
	Content    []map[string]string `bson:"content" json:"content" map:"content"`             //key=>相关集合的字段，values=>相关集合的内容
	Status     int                 `bson:"status" json:"status" map:"status"`                //只有状态为predefined.AuditStatusPass时，才会真正存入相关集合的相关字段
	CreateTime time.Time           `bson:"create_time" json:"create_time" map:"create_time"` //创建时间
	UpdateTime time.Time           `bson:"update_time" json:"update_time" map:"update_time"` //更新时间
}
