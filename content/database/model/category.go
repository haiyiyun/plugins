package model

import (
	"time"

	"github.com/haiyiyun/mongodb/geometry"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID                          primitive.ObjectID   `json:"_id" bson:"_id,omitempty" map:"_id,omitempty"`
	ParentID                    primitive.ObjectID   `json:"parent_id" bson:"parent_id" map:"parent_id"`                                                                         //父级ID
	Type                        int                  `json:"type" bson:"type" map:"type"`                                                                                        //分类类型
	Name                        string               `json:"name" bson:"name" map:"name"`                                                                                        //分类名
	Tags                        []string             `json:"tags" bson:"tags" map:"tags"`                                                                                        //标签
	LimitUserAtLeastLevel       int                  `json:"limit_user_at_least_level" bson:"limit_user_at_least_level" map:"limit_user_at_least_level"`                         //限制用户至少级别
	OnlyUserIDNotLimitUserLevel []primitive.ObjectID `json:"only_user_id_not_limit_user_level" bson:"only_user_id_not_limit_user_level" map:"only_user_id_not_limit_user_level"` //只有哪些user_id不受等级限制
	LimitUserRole               []string             `json:"limit_user_role" bson:"limit_user_role" map:"limit_user_role"`                                                       //限制用户角色
	OnlyUserIDNotLimitUserRole  []primitive.ObjectID `json:"only_user_id_not_limit_user_role" bson:"only_user_id_not_limit_user_role" map:"only_user_id_not_limit_user_role"`    //只有哪些user_id不受角色限制
	LimitUserTag                []string             `json:"limit_user_tag" bson:"limit_user_tag" map:"limit_user_tag"`                                                          //限制用户TAG
	OnlyUserIDNotLimitUserTag   []primitive.ObjectID `json:"only_user_id_not_limit_user_tag" bson:"only_user_id_not_limit_user_tag" map:"only_user_id_not_limit_user_tag"`       //只有哪些user_id不受tag限制
	Location                    geometry.Point       `json:"location" bson:"location,omitempty" map:"location,omitempty"`
	Visibility                  int                  `json:"visibility" bson:"visibility" map:"visibility"` //可见度
	ExtraData                   string               `json:"extra_data" bson:"extra_data" map:"extra_data"` //额外扩展信息数据，可灵活使用，比如将相关额外信息json后存入
	Status                      int                  `json:"status" bson:"status" map:"status"`
	Enable                      bool                 `json:"enable" bson:"enable" map:"enable"`
	CreateTime                  time.Time            `json:"create_time" bson:"create_time" map:"create_time"`
	UpdateTime                  time.Time            `json:"update_time" bson:"update_time" map:"update_time"`
}
