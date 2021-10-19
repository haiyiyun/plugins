package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Share struct {
	ID                 primitive.ObjectID `json:"_id" bson:"_id,omitempty" map:"_id,omitempty"`
	Type               int                `json:"type" bson:"type" map:"type"` //分享类型
	UserID             primitive.ObjectID `json:"user_id" bson:"user_id" map:"user_id"`
	SharerUserID       primitive.ObjectID `json:"sharer_user_id" bson:"sharer_user_id" map:"sharer_user_id"`                   //分享者user_id
	PlatformCategoryID primitive.ObjectID `json:"platform_category_id" bson:"platform_category_id" map:"platform_category_id"` //平台分类ID
	ObjectID           primitive.ObjectID `json:"object_id" bson:"object_id" map:"object_id"`
	CreateTime         time.Time          `json:"create_time" bson:"create_time" map:"create_time"`
	UpdateTime         time.Time          `json:"update_time" bson:"update_time" map:"update_time"`
}
