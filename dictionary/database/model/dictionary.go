package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DictionaryValue struct {
	Lable  string `bson:"lable" json:"lable" map:"lable"` //标签
	Key    string `bson:"key" json:"key" map:"key"`       //键名
	Value  int    `bson:"value" json:"value" map:"value"` //键值
	Order  int    `bson:"order" json:"order" map:"order"` //排序
	Enable bool   `bson:"enable" json:"enable" map:"enable"`
}

type Dictionary struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id" map:"_id"`
	Name       string             `bson:"name" json:"name" map:"name"`
	Key        string             `bson:"key" json:"key" map:"key"`
	Values     []DictionaryValue  `bson:"values" json:"values" map:"values"`
	Remark     string             `bson:"remark" json:"remark" map:"remark"`
	Enable     bool               `bson:"enable" json:"enable" map:"enable"`
	Delete     bool               `bson:"delete" json:"delete" map:"delete"`
	CreateTime time.Time          `bson:"create_time" json:"create_time" map:"create_time"` //创建时间
	UpdateTime time.Time          `bson:"update_time" json:"update_time" map:"update_time"` //更新时间
}
