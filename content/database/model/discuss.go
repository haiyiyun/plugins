package model

import (
	"time"

	"github.com/haiyiyun/mongodb/geometry"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Discuss struct {
	ID             primitive.ObjectID   `json:"_id" bson:"_id,omitempty" map:"_id,omitempty"`
	Type           int                  `json:"type" bson:"type" map:"type"`
	ObjectID       primitive.ObjectID   `json:"object_id" bson:"object_id" map:"object_id"`                      //根据type的类型，确定其类型的_id作为object_id
	PublishUserID  primitive.ObjectID   `json:"publish_user_id" bson:"publish_user_id" map:"publish_user_id"`    //发布者user_id，为空代表系统发布
	AtUsers        []primitive.ObjectID `json:"at_users" bson:"at_users" map:"at_users"`                         //@用户user_id
	ReplyDiscussID primitive.ObjectID   `json:"reply_discuss_id" bson:"reply_discuss_id" map:"reply_discuss_id"` //回复评论id
	Text           string               `json:"text" bson:"text" map:"text"`                                     //纯文本
	Location       geometry.Point       `json:"location" bson:"location,omitempty" map:"location,omitempty"`
	LikedUser      []primitive.ObjectID `json:"liked_user" bson:"liked_user" map:"liked_user"` //喜欢用户user_id
	HatedUser      []primitive.ObjectID `json:"hated_user" bson:"hated_user" map:"hated_user"` //讨厌用户user_id
	Visibility     int                  `json:"visibility" bson:"visibility" map:"visibility"` //可见度
	Status         int                  `json:"status" bson:"status" map:"status"`
	CreateTime     time.Time            `json:"create_time" bson:"create_time" map:"create_time"`
	UpdateTime     time.Time            `json:"update_time" bson:"update_time" map:"update_time"`
}
