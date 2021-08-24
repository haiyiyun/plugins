package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenSignInfo struct {
	IP        string     `bson:"ip" json:"ip" map:"ip"`
	UserAgent string     `bson:"user_agent" json:"user_agent" map:"user_agent"`
	Geo       [2]float64 `bson:"geo" json:"geo" map:"geo"`
}

type Token struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id" map:"_id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id" map:"user_id"`
	UserName    string             `bson:"user_name" json:"user_name" map:"user_name"`
	TokenType   string             `bson:"token_type" json:"token_type" map:"token_type"` //类型
	Token       string             `bson:"token" json:"token" map:"token"`
	SignInfo    TokenSignInfo      `bson:"sign_info" json:"sign_info" map:"sign_info"`
	ExpiredTime time.Time          `bson:"expired_time" json:"expired_time" map:"expired_time"` //过期时间
	CreateTime  time.Time          `bson:"create_time" json:"create_time" map:"create_time"`    //创建时间
	UpdateTime  time.Time          `bson:"update_time" json:"update_time" map:"update_time"`    //更新时间
}
