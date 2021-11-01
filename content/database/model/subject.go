package model

import (
	"time"

	"github.com/haiyiyun/mongodb/geometry"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subject struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id,omitempty" map:"_id,omitempty"`
	PublishUserID primitive.ObjectID `json:"publish_user_id" bson:"publish_user_id" map:"publish_user_id"` //发布者，为空代表系统发布
	Type          int                `json:"type" bson:"type" map:"type"`                                  //主题类型
	Subject       string             `json:"subject" bson:"subject" map:"subject"`                         //主题
	Cover         string             `json:"cover" bson:"cover" map:"cover"`                               //封面图片
	Description   string             `json:"description" bson:"description" map:"description"`             //描述
	Text          string             `json:"text" bson:"text" map:"text"`                                  //文本
	Images        []string           `json:"images" bson:"images" map:"images"`                            //图片
	Video         string             `json:"video" bson:"video" map:"video"`                               //视频
	Voice         string             `json:"voice" bson:"voice" map:"voice"`                               //语音
	UserTags      []string           `json:"user_tags" bson:"user_tags" map:"user_tags"`                   //用户可编辑的标签
	Tags          []string           `json:"tags" bson:"tags" map:"tags"`                                  //标签，包括用户标签
	Location      geometry.Point     `json:"location" bson:"location,omitempty" map:"location,omitempty"`
	Visibility    int                `json:"visibility" bson:"visibility" map:"visibility"` //可见度
	Status        int                `json:"status" bson:"status" map:"status"`
	Enable        bool               `json:"enable" bson:"enable" map:"enable"`
	CreateTime    time.Time          `json:"create_time" bson:"create_time" map:"create_time"`
	UpdateTime    time.Time          `json:"update_time" bson:"update_time" map:"update_time"`
}
