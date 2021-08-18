package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Log struct {
	ID             primitive.ObjectID  `bson:"_id,omitempty" json:"_id" map:"_id"`
	Type           string              `bson:"type" json:"type" map:"type"`
	UserID         primitive.ObjectID  `bson:"user_id" json:"user_id" map:"user_id"`
	User           string              `bson:"user" json:"user" map:"user"`
	Method         string              `bson:"method" json:"method" map:"method"`
	Referer        string              `bson:"referer" json:"referer" map:"referer"`
	Path           string              `bson:"path" json:"path" map:"path"`
	Query          string              `bson:"query" json:"query" map:"query"`
	IP             string              `bson:"ip" json:"ip" map:"ip"`
	RequestHeader  map[string][]string `bson:"request_header" json:"request_header" map:"request_header"`
	RequestPayload string              `bson:"request_payload" json:"request_payload" map:"request_payload"`
	ResponseHeader map[string][]string `bson:"response_header" json:"response_header" map:"response_header"`
	ResponseData   string              `bson:"response_data" json:"response_data" map:"response_data"`
	DeleteTime     time.Time           `json:"delete_time" bson:"delete_time" map:"delete_time"` //自动删除时间
	CreateTime     time.Time           `bson:"create_time" json:"create_time" map:"create_time"` //创建时间
	UpdateTime     time.Time           `bson:"update_time" json:"update_time" map:"update_time"` //更新时间
}
