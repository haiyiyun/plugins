package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageLink struct {
	Type   int    `json:"type" bson:"type" map:"type"`       //链接类型
	URL    string `json:"url" bson:"url" map:"url"`          //链接地址
	Normal bool   `json:"normal" bson:"normal" map:"normal"` //是否正常链接
	Inline bool   `json:"inline" bson:"inline" map:"inline"` //是否内嵌，还是跳转
	Block  bool   `json:"block" bson:"block" map:"block"`    //是否需要屏蔽跳转
}

type Message struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id,omitempty" map:"_id,omitempty"`
	Type       int                `json:"type" bson:"type" map:"type"`                         //消息类型
	FromUserID primitive.ObjectID `json:"from_user_id" bson:"from_user_id" map:"from_user_id"` //系统消息时，from_user_id为NilObjectID
	ToUserID   primitive.ObjectID `json:"to_user_id" bson:"to_user_id" map:"to_user_id"`
	Text       string             `json:"text" bson:"text" map:"text"`
	RichText   string             `json:"rich_text" bson:"rich_text" map:"rich_text"` //富文本，如：支持链接，排版等
	Image      string             `json:"image" bson:"image" map:"image"`
	Voice      string             `json:"voice" bson:"voice" map:"voice"`
	Video      string             `json:"video" bson:"video" map:"video"`
	File       string             `json:"file" bson:"file" map:"file"`
	Link       MessageLink        `json:"link" bson:"link" map:"link"`
	Status     int                `json:"status" bson:"status" map:"status"`
	CreateTime time.Time          `json:"create_time" bson:"create_time" map:"create_time"`
	UpdateTime time.Time          `json:"update_time" bson:"update_time" map:"update_time"`
}
