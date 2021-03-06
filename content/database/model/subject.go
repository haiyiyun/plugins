package model

import (
	"time"

	"github.com/haiyiyun/mongodb/geometry"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subject struct {
	ID                          primitive.ObjectID   `json:"_id" bson:"_id,omitempty" map:"_id,omitempty"`
	PublishUserID               primitive.ObjectID   `json:"publish_user_id" bson:"publish_user_id" map:"publish_user_id"` //发布者，为空代表系统发布
	Type                        int                  `json:"type" bson:"type" map:"type"`                                  //主题类型
	Subject                     string               `json:"subject" bson:"subject" map:"subject"`                         //主题
	Cover                       string               `json:"cover" bson:"cover" map:"cover"`                               //封面图片
	Description                 string               `json:"description" bson:"description" map:"description"`             //描述
	Text                        string               `json:"text" bson:"text" map:"text"`                                  //文本
	Images                      []string             `json:"images" bson:"images" map:"images"`                            //图片
	Video                       string               `json:"video" bson:"video" map:"video"`                               //视频
	Voice                       string               `json:"voice" bson:"voice" map:"voice"`                               //语音
	UserTags                    []string             `json:"user_tags" bson:"user_tags" map:"user_tags"`                   //用户可编辑的标签
	Tags                        []string             `json:"tags" bson:"tags" map:"tags"`                                  //标签，包括用户标签
	Location                    geometry.Point       `json:"location" bson:"location,omitempty" map:"location,omitempty"`
	LimitUserAtLeastLevel       int                  `json:"limit_user_at_least_level" bson:"limit_user_at_least_level" map:"limit_user_at_least_level"`                         //限制用户至少级别
	OnlyUserIDNotLimitUserLevel []primitive.ObjectID `json:"only_user_id_not_limit_user_level" bson:"only_user_id_not_limit_user_level" map:"only_user_id_not_limit_user_level"` //只有哪些user_id不受等级限制
	LimitUserRole               []string             `json:"limit_user_role" bson:"limit_user_role" map:"limit_user_role"`                                                       //限制用户角色
	OnlyUserIDNotLimitUserRole  []primitive.ObjectID `json:"only_user_id_not_limit_user_role" bson:"only_user_id_not_limit_user_role" map:"only_user_id_not_limit_user_role"`    //只有哪些user_id不受角色限制
	LimitUserTag                []string             `json:"limit_user_tag" bson:"limit_user_tag" map:"limit_user_tag"`
	OnlyUserIDNotLimitUserTag   []primitive.ObjectID `json:"only_user_id_not_limit_user_tag" bson:"only_user_id_not_limit_user_tag" map:"only_user_id_not_limit_user_tag"` //只有哪些user_id不受tag限制
	Visibility                  int                  `json:"visibility" bson:"visibility" map:"visibility"`                                                                //可见度
	ExtraData                   string               `json:"extra_data" bson:"extra_data" map:"extra_data"`                                                                //额外扩展信息数据，可灵活使用，比如将相关额外信息json后存入
	Status                      int                  `json:"status" bson:"status" map:"status"`
	Enable                      bool                 `json:"enable" bson:"enable" map:"enable"`
	CreateTime                  time.Time            `json:"create_time" bson:"create_time" map:"create_time"`
	UpdateTime                  time.Time            `json:"update_time" bson:"update_time" map:"update_time"`
}
