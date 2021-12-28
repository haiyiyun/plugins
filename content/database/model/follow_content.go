package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//关注内容
//当内容发布时，会在这个集合里插入一条记录，关注者关联此记录的UserID可获取被关注的ContentID,通过ContentID获取具体内容信息
//当取消关注时，通过FollowRelationshipID会清除相关FollowContent的记录
//是否被阅读过，阅读过的内容，可选择是否要在一段时间后自动删除,比如，一般是：一个月后
type FollowContent struct {
	ID                   primitive.ObjectID `json:"_id" bson:"_id,omitempty" map:"_id,omitempty"`
	FollowRelationshipID primitive.ObjectID `json:"follow_relationship_id" bson:"follow_relationship_id" map:"follow_relationship_id"` //关注关系ID
	Type                 int                `json:"type" bson:"type" map:"type"`                                                       //关注内容类型，与内容发布类型对应
	UserID               primitive.ObjectID `json:"user_id" bson:"user_id" map:"user_id"`                                              //关注者user_id
	ContentID            primitive.ObjectID `json:"content_id" bson:"content_id" map:"content_id"`                                     //被关注对象的内容的ID，如动态的ID,评论的ID
	ReadedTime           time.Time          `json:"readed_time" bson:"readed_time" map:"readed_time"`
	ExtensionID          primitive.ObjectID `bson:"extension_id" json:"extension_id" map:"extension_id"` //扩展ID,冗余字段，从follow_relationship表硬传过来，用来检索用
	CreateTime           time.Time          `json:"create_time" bson:"create_time" map:"create_time"`
	UpdateTime           time.Time          `json:"update_time" bson:"update_time" map:"update_time"`
}
