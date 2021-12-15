package predefined

import "go.mongodb.org/mongo-driver/bson/primitive"

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

type RequestServeTags struct {
	Tags []string `form:"tags,omitempty"`
}

type RequestServeUserTags struct {
	UserTags []string `form:"user_tags,omitempty"`
}

type RequestServeAtUsers struct {
	AtUsers []primitive.ObjectID `form:"at_users,omitempty"`
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
	Subject string `form:"subject" validate:"required"`
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
	RequestServeObjectID
	RequestServeDiscussType
	RequestServeAtUsers
	ReplyDiscussID primitive.ObjectID `form:"reply_discuss_id,omitempty"`
	RequestServeDiscussVisibility
	RequestServeLongitudeLatitude
	Text string `form:"text" validate:"required"`
}

type RequestServeDiscussList struct {
	RequestServeObjectID
	RequestServeDiscussType
	RequestServeDiscussVisibility
	RequestServeLongitudeLatitude
	RequestServeDistance
	RequestServePublishUserID
}

type RequestServeContentCreateType struct {
	Type int `form:"type" validate:"oneof=0 1 2 3 4 5 6"`
}

type RequestServeContentListType struct {
	Type int `form:"type" validate:"oneof=0 1 2 3 4 5 6 1000 1001 1002 1003 1004 1005 1006"`
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
	AssociateID primitive.ObjectID `form:"associate_id" validate:"required_unless=type 0"`
	CategoryID  primitive.ObjectID `form:"category_id"`
	SubjectID   primitive.ObjectID `form:"subject_id"`
	RequestServeAtUsers
	Author      string   `form:"author"`
	Title       string   `form:"title" validate:"required"`
	Cover       string   `form:"cover" validate:"required"`
	Description string   `form:"description,omitempty"`
	Video       string   `form:"video" validate:"required_with=type 0 type 1"`
	Voice       string   `form:"voice" validate:"required_with=type 2 type 3"`
	Images      []string `form:"images" validate:"required_with=type 4 type 5"`
	Content     string   `form:"content" validate:"required_with=type 4 type 6"`
	ContentLink string   `form:"content_link"`
	RequestServeLongitudeLatitude
	RequestServeUserTags
	RequestServeContentVisibility
	Copy           bool `form:"copy"`
	ForbidForward  bool `form:"forbid_forward"`
	ForbidDownload bool `form:"forbid_download"`
	ForbidDiscuss  bool `form:"forbid_discuss"`
}

type RequestServeContentList struct {
	RequestServeContentListType
	RequestServeContentPublishType
	RequestServeContentVisibility
	RequestServeTags
	RequestServePublishUserID
	CategoryID primitive.ObjectID `form:"category_id,omitempty"`
	SubjectID  primitive.ObjectID `form:"subject_id,omitempty"`
	RequestServeLongitudeLatitude
	RequestServeDistance
}

type RequestServeContentDetail struct {
	RequestServeID
}
