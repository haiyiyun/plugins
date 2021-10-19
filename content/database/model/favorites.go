package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//收藏夹
type Favorites struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id,omitempty" map:"_id,omitempty"`
	Type       int                `json:"type" bson:"type" map:"type"`                //收藏的类型，如：动态，文章
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id" map:"user_id"`       //关注者user_id
	ObjectID   primitive.ObjectID `json:"object_id" bson:"object_id" map:"object_id"` //被收藏对象的ID，如动态的ID
	Content    string             `json:"content" bson:"content" map:"content"`       //收藏内容，备份使用，因为有的时候ObjectID可能被删除
	Status     int                `json:"status" bson:"status" map:"status"`
	CreateTime time.Time          `json:"create_time" bson:"create_time" map:"create_time"`
	UpdateTime time.Time          `json:"update_time" bson:"update_time" map:"update_time"`
}
