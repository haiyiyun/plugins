package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContactsApplyReply struct {
	Applyer   bool      `bson:"applyer" json:"applyer" map:"applyer"`          //是否申请者
	Message   string    `bson:"message" json:"message" map:"message"`          //回复信息
	ReplyTime time.Time `bson:"reply_time" json:"reply_time" map:"reply_time"` //回复时间
}

type ContactsApply struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty" json:"_id" map:"_id"`
	UserID        primitive.ObjectID   `bson:"user_id" json:"user_id" map:"user_id"`
	ApplyerUserID primitive.ObjectID   `bson:"applyer_user_id" json:"applyer_user_id" map:"applyer_user_id"`
	ApplyReason   string               `bson:"apply_reason" json:"apply_reason" map:"apply_reason"` //申请理由
	Replys        []ContactsApplyReply `bson:"replys" json:"replys" map:"replys"`
	Refuse        bool                 `bson:"refuse" json:"refuse" map:"refuse"`                      //是否拒绝，Refuse后，一段时间会设置Delete为true
	RefuseReason  string               `bson:"refuse_reason" json:"refuse_reason" map:"refuse_reason"` //拒绝理由
	RefuseTime    time.Time            `bson:"refuse_time" json:"refuse_time" map:"refuse_time"`
	Pass          bool                 `bson:"pass" json:"pass" map:"pass"` //是否通过
	PassTime      time.Time            `bson:"pass_time" json:"pass_time" map:"pass_time"`
	Delete        bool                 `bson:"delete" json:"delete" map:"delete"`                //Pass后，就会设置Delete为true；Refuse后，一段时间会设置Delete为true，比如：2天；或者UpdateTime时间后会自动设为true，比如：3天
	CreateTime    time.Time            `bson:"create_time" json:"create_time" map:"create_time"` //创建时间
	UpdateTime    time.Time            `bson:"update_time" json:"update_time" map:"update_time"` //更新时间
}
