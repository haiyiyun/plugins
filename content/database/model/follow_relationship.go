package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//关注关系
//互相确认成为联系人时，也算建立了关注关系，并且算互相关注
type FollowRelationship struct {
	ID                primitive.ObjectID `json:"_id" bson:"_id,omitempty" map:"_id,omitempty"`
	Type              int                `json:"type" bson:"type" map:"type"`                                                 //关注关系类型
	UserID            primitive.ObjectID `json:"user_id" bson:"user_id" map:"user_id"`                                        //关注者user_id
	ObjectID          primitive.ObjectID `json:"object_id" bson:"object_id" map:"object_id"`                                  //被关注对象的ID
	ObjectOwnerUserID primitive.ObjectID `bson:"object_owner_user_id" json:"object_owner_user_id" map:"object_owner_user_id"` //被关注对象的拥有者的user_id
	Mutual            bool               `json:"mutual" bson:"mutual" map:"mutual"`                                           //是否互相关注对方
	Stealth           bool               `bson:"stealth" json:"stealth" map:"stealth"`                                        //是否隐身关注
	ExtensionID       primitive.ObjectID `bson:"extension_id" json:"extension_id" map:"extension_id"`                         //扩展ID
	CreateTime        time.Time          `json:"create_time" bson:"create_time" map:"create_time"`
	UpdateTime        time.Time          `json:"update_time" bson:"update_time" map:"update_time"`
}
