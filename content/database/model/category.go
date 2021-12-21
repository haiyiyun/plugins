package model

import (
	"time"

	"github.com/haiyiyun/mongodb/geometry"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id,omitempty" map:"_id,omitempty"`
	ParentID   primitive.ObjectID `json:"parent_id" bson:"parent_id" map:"parent_id"` //父级ID
	Type       int                `json:"type" bson:"type" map:"type"`                //分类类型
	Name       string             `json:"name" bson:"name" map:"name"`                //分类名
	Tags       []string           `json:"tags" bson:"tags" map:"tags"`                //标签
	Location   geometry.Point     `json:"location" bson:"location,omitempty" map:"location,omitempty"`
	Visibility int                `json:"visibility" bson:"visibility" map:"visibility"` //可见度
	ExtraData  string             `json:"extra_data" bson:"extra_data" map:"extra_data"` //额外扩展信息数据，可灵活使用，比如将相关额外信息json后存入
	Status     int                `json:"status" bson:"status" map:"status"`
	Enable     bool               `json:"enable" bson:"enable" map:"enable"`
	CreateTime time.Time          `json:"create_time" bson:"create_time" map:"create_time"`
	UpdateTime time.Time          `json:"update_time" bson:"update_time" map:"update_time"`
}
