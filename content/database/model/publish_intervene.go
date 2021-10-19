package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//发布干预，用来控制发布内容时的状态。出现在此集合中的用户，发布时会被进行状态处理，不能直接正常发布
type PublishIntervene struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id,omitempty" map:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id" map:"user_id"` //user_id为primitive.NilObjectID时，代表所有用户
	Type       int                `json:"type" bson:"type" map:"type"`
	Reason     string             `json:"reason" bson:"reason" map:"reason"` //干预原因
	Status     int                `json:"status" bson:"status" map:"status"`
	StartTime  time.Time          `json:"start_time" bson:"start_time" map:"start_time"` //干预开始时间
	EndTime    time.Time          `json:"end_time" bson:"end_time" map:"end_time"`       //干预结束时间
	CreateTime time.Time          `json:"create_time" bson:"create_time" map:"create_time"`
	UpdateTime time.Time          `json:"update_time" bson:"update_time" map:"update_time"`
}
