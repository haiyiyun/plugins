package predefined

import (
	"github.com/haiyiyun/utils/help"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RequestServeLongitudeLatitude struct {
	Longitude float64 `form:"longitude,omitempty"` //经度
	Latitude  float64 `form:"latitude,omitempty"`  //维度
}

type RequestServeDistance struct {
	MaxDistance float64 `form:"max_distance,omitempty" validate:"numeric"`
	MinDistance float64 `form:"min_distance,omitempty" validate:"numeric"`
}

type RequestServeParentID struct {
	ParentID primitive.ObjectID `form:"parent_id,omitempty"`
}

type RequestServePublishUserID struct {
	PublishUserID primitive.ObjectID `form:"publish_user_id,omitempty"`
}

type RequestServeID struct {
	ID primitive.ObjectID `form:"id" validate:"required"`
}

type RequestServeObjectID struct {
	ObjectID primitive.ObjectID `form:"object_id,omitempty"`
}

type RequestServeObjectIDRequired struct {
	ObjectID primitive.ObjectID `form:"object_id" validate:"required"`
}

type RequestServeTags struct {
	Tags []string `form:"tags"`
}

type RequestServeUserTags struct {
	UserTags []string `form:"user_tags"`
}

type RequestServeAtUsers struct {
	AtUsers []primitive.ObjectID `form:"at_users"`
}

type RequestServeCategorySystemType struct {
	Type int `form:"type" validate:"oneof=0 1 2 3 4"`
}

type RequestServeCategoryVisibility struct {
	Visibility int `form:"visibility" validate:"oneof=1 2 3 4 5 6 7 8 9"`
}

type RequestServeCategoryList struct {
	RequestServeParentID
	RequestServeCategorySystemType
	RequestServeCategoryVisibility
	RequestServeLongitudeLatitude
	RequestServeDistance
	RequestServeTags
}

type RequestServeSubjectCreateType struct {
	Type int `form:"type" validate:"oneof=0 1"`
}

type RequestServeSubjectListType struct {
	Type int `form:"type" validate:"oneof=0 1 1000 1001"`
}

type RequestServeSubjectVisibility struct {
	Visibility int `form:"visibility" validate:"oneof=0 1 2 3 4 5 6 7 8 9"`
}

type RequestServeSubjectCreate struct {
	RequestServeSubjectCreateType
	RequestServeSubjectVisibility
	RequestServeLongitudeLatitude
	RequestServeUserTags
	ExtraData string `form:"extra_data,omitempty"`
	Subject   string `form:"subject" validate:"required"`
}

type RequestServeSubjectList struct {
	RequestServePublishUserID
	RequestServeSubjectListType
	RequestServeSubjectVisibility
	RequestServeLongitudeLatitude
	RequestServeDistance
	RequestServeTags
}

type RequestServeDiscussType struct {
	Type int `form:"type" validate:"oneof=0 1 2 3"`
}

type RequestServeDiscussVisibility struct {
	Visibility int `form:"visibility" validate:"oneof=0 1 2 3 4 5 6 7 8 9"`
}

type RequestServeDiscussCreate struct {
	RequestServeObjectIDRequired
	RequestServeDiscussType
	RequestServeAtUsers
	ReplyDiscussID primitive.ObjectID `form:"reply_discuss_id,omitempty"`
	RequestServeDiscussVisibility
	RequestServeLongitudeLatitude
	Evaluation float64 `form:"evaluation"`
	Text       string  `form:"text" validate:"required"`
}

type RequestServeDiscussList struct {
	Types []int `form:"types,omitempty" validate:"required_with=ObjectID,gte=0,dive,oneof=0 1 2 3"`
	RequestServeObjectID
	RequestServeDiscussVisibility
	RequestServePublishUserID
	RequestServeLongitudeLatitude
	RequestServeDistance
}

type RequestServeContentCreateType struct {
	Type int `form:"type" validate:"oneof=0 1 2 3 4 5 6"`
}

type RequestServeContentListType struct {
	Type int `form:"type" validate:"oneof=0 1 2 3 4 5 6 1000 1001 1002 1003 1004 1005 1006"`
}

type RequestServeContentTypes struct {
	Types []int `form:"types,omitempty" validate:"gte=0,dive,oneof=0 1 2 3 4 5 6 1000 1001 1002 1003 1004 1005 1006"`
}

type RequestServeContentPublishType struct {
	PublishType int `form:"publish_type" validate:"oneof=0 1 2 3"`
}

type RequestServeContentAssociateType struct {
	AssociateType int `form:"associate_type" validate:"oneof=0 1 2 3 4 5 6"`
}

type RequestServeContentVisibility struct {
	Visibility int `form:"visibility" validate:"oneof=0 1 2 3 4 5 6 7 8 9"`
}

type RequestServeContentCreate struct {
	RequestServeContentCreateType
	RequestServeContentPublishType
	RequestServeContentAssociateType
	AssociateID        primitive.ObjectID `form:"associate_id" validate:"required_unless=AssociateType 0"`
	LimitAssociateType int                `form:"limit_associate_type" validate:"oneof=0 1 2 3 4 5 6"`
	LimitAssociateNum  int                `form:"limit_associate_num"`
	CategoryID         primitive.ObjectID `form:"category_id"`
	SubjectID          primitive.ObjectID `form:"subject_id"`
	RequestServeAtUsers
	Author      string   `form:"author"`
	Title       string   `form:"title" validate:"required"`
	Cover       string   `form:"cover,omitempty"`
	Description string   `form:"description,omitempty"`
	Video       string   `form:"video" validate:"required_if=Type 0,required_if=Type 1"`
	Voice       string   `form:"voice" validate:"required_if=Type 2,required_if=Type 3"`
	Images      []string `form:"images" validate:"required_if=Type 4,required_if=Type 5"`
	Content     string   `form:"content" validate:"required_if=Type 4,required_if=Type 6"`
	ContentLink string   `form:"content_link"`
	RequestServeLongitudeLatitude
	RequestServeUserTags
	RequestServeContentVisibility
	Value                                  int                  `form:"value"`
	HideDetail                             bool                 `form:"hide_detail"`
	OnlyUserIDShowDetail                   []primitive.ObjectID `form:"only_user_id_show_detail"`
	Copy                                   bool                 `form:"copy"`
	LimitUserAtLeastLevel                  int                  `form:"limit_user_at_least_level"`
	OnlyUserIDNotLimitUserLevel            []primitive.ObjectID `form:"only_user_id_not_limit_user_level"`
	LimitUserRole                          []string             `form:"limit_user_role"`
	OnlyUserIDNotLimitUserRole             []primitive.ObjectID `form:"only_user_id_not_limit_user_role"`
	LimitUserTag                           []string             `form:"limit_user_tag"`
	OnlyUserIDNotLimitUserTag              []primitive.ObjectID `form:"only_user_id_not_limit_user_tag"`
	OnlyUserIDDiscuss                      []primitive.ObjectID `form:"only_user_id_discuss"`
	OnlyUserIDCanReplyDiscuss              []primitive.ObjectID `form:"only_publish_user_id_can_reply_discuss"`
	OnlyUserIDCanNotReplyDiscuss           []primitive.ObjectID `form:"only_publish_user_id_can_not_reply_discuss"`
	LimitAllDiscussNum                     int                  `form:"limit_all_discuss_num"`
	LimitPublishUserDiscussNum             int                  `form:"limit_publish_user_discuss_num"`
	LimitNotPublishUserAllUserDiscussNum   int                  `form:"limit_not_publish_user_all_user_discuss_num"`
	LimitNotPublishUserEveryUserDiscussNum int                  `form:"limit_not_publish_user_every_user_discuss_num"`
	HideDiscuss                            bool                 `form:"hide_discuss"`
	OnlyUserIDShowDiscuss                  []primitive.ObjectID `form:"only_user_id_show_discuss"`
	ForbidForward                          bool                 `form:"forbid_forward"`
	ForbidDownload                         bool                 `form:"forbid_download"`
	ForbidDiscuss                          bool                 `form:"forbid_discuss"`
	StartTime                              help.DateTime        `form:"start_time"`
	EndTime                                help.DateTime        `form:"end_time"`
	ExtraData                              string               `form:"extra_data"`
}

type RequestServeContentList struct {
	RequestServeContentTypes
	RequestServeContentPublishType
	RequestServeContentVisibility
	RequestServeTags
	RequestServePublishUserID
	CategoryID primitive.ObjectID `form:"category_id,omitempty"`
	SubjectID  primitive.ObjectID `form:"subject_id,omitempty"`
	RequestServeLongitudeLatitude
	RequestServeDistance
	LimitUserAtLeastLevel int           `form:"limit_user_at_least_level,omitempty"`
	LimitUserRole         []string      `form:"limit_user_role,omitempty"`
	LimitUserTag          []string      `form:"limit_user_tag,omitempty"`
	ValueZero             bool          `form:"value_zero,omitempty"`
	ValueGte              int           `form:"value_gte,omitempty"`
	ValueLte              int           `form:"value_lte,omitempty"`
	DiscussTotalZero      bool          `form:"discuss_total_zero,omitempty"`
	DiscussTotalGte       int           `form:"discuss_total_gte,omitempty"`
	DiscussTotalLte       int           `form:"discuss_total_lte,omitempty"`
	InTime                bool          `form:"in_time,omitempty"`
	StartTime             help.DateTime `form:"start_time,omitempty"`
	EndTime               help.DateTime `form:"end_time,omitempty"`
}

type RequestServeContentDetail struct {
	RequestServeID
}
