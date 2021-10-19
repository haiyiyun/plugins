package model

import (
	"time"

	"github.com/haiyiyun/mongodb/geometry"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Discuss struct {
	ID             primitive.ObjectID   `json:"_id" bson:"_id,omitempty" map:"_id,omitempty"`
	Type           int                  `json:"type" bson:"type" map:"type"`
	ObjectID       primitive.ObjectID   `json:"object_id" bson:"object_id" map:"object_id"`
	UserID         primitive.ObjectID   `json:"user_id" bson:"user_id" map:"user_id"`
	AtUser         []primitive.ObjectID `json:"at_user" bson:"at_user" map:"at_user"`                            //@用户user_id
	ReplyDiscussID primitive.ObjectID   `json:"reply_discuss_id" bson:"reply_discuss_id" map:"reply_discuss_id"` //回复评论id
	Text           string               `json:"text" bson:"text" map:"text"`
	Location       geometry.Point       `json:"location" bson:"location,omitempty" map:"location,omitempty"`
	LikedUser      []primitive.ObjectID `json:"liked_user" bson:"liked_user" map:"liked_user"` //喜欢用户user_id
	HatedUser      []primitive.ObjectID `json:"hated_user" bson:"hated_user" map:"hated_user"` //讨厌用户user_id
	Visibility     int                  `json:"visibility" bson:"visibility" map:"visibility"` //可见度
	Status         int                  `json:"status" bson:"status" map:"status"`
	CreateTime     time.Time            `json:"create_time" bson:"create_time" map:"create_time"`
	UpdateTime     time.Time            `json:"update_time" bson:"update_time" map:"update_time"`
}
