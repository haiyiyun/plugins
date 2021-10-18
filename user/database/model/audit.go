package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//审核，需要经过审核的，在此选择类型并添加需审核的字段，
//相关集合写入前会被拦截，并添加到相应的审核记录集合里面进行审核，
//审核通过了，才能正式入库
type Audit struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id" map:"_id"`
	Type       int                `bson:"type" json:"type" map:"type"`                      //类型，如profile的头像等
	Fields     []string           `bson:"fields" json:"fields" map:"fields"`                //需要审核的字段
	Enable     bool               `bson:"enable" json:"enable" map:"enable"`                //是否启用此审核条件
	CreateTime time.Time          `bson:"create_time" json:"create_time" map:"create_time"` //创建时间
	UpdateTime time.Time          `bson:"update_time" json:"update_time" map:"update_time"` //更新时间
}
