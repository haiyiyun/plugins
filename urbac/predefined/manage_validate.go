package predefined

import "go.mongodb.org/mongo-driver/bson/primitive"

type RequestManageApplicationUpdate struct {
	ParentPath   string `form:"parent_path"`
	Type         string `form:"type" validate:"required"`
	Name         string `form:"name" validate:"required"`
	Path         string `form:"path" validate:"required"`
	Level        string `form:"level" validate:"required"`
	Order        int    `form:"order" validate:"numeric"`
	Enable       bool   `form:"enable"`
	MetaHideMenu bool   `form:"meta_hide_menu"`
	MetaTitle    string `form:"meta_title" validate:"required"`
	MetaIcon     string `form:"meta_icon"`
	MetaFrameSrc string `form:"meta_frame_src"`
}

type RequestManageNamePath struct {
	Name string `form:"name" validate:"required"`
	Path string `form:"path" validate:"required"`
}

type RequestManageLevelPath struct {
	Level string `form:"level" validate:"required"`
	Path  string `form:"path" validate:"required"`
}

type RequestManageEnable struct {
	Enable bool `form:"enable"`
}

type RequestManageHide struct {
	Hide bool `form:"hide"`
}

type RequestManageApplicationEnable struct {
	RequestManageLevelPath
	RequestManageEnable
}

type RequestManageApplicationHide struct {
	RequestManageLevelPath
	RequestManageHide
}

type RequestManageUsername struct {
	Username string `form:"username" validate:"required"`
}

type RequestManagePassword struct {
	Password string `form:"password" validate:"required"`
}

type RequestManageLongitudeLatitude struct {
	Longitude float64 `form:"longitude,omitempty"` //经度
	Latitude  float64 `form:"latitude,omitempty"`  //维度
}

type RequestManageLogin struct {
	RequestManageUsername
	RequestManagePassword
	RequestManageLongitudeLatitude
}

type RequestManageRefresh struct {
	RequestManageLongitudeLatitude
}

type RequestManageRole struct {
	Name  string `form:"name" validate:"required"`
	Scope int    `form:"scope" validate:"gte=0"`
	RequestManageEnable
	Applications string `form:"applications" validate:"json"`
}

type RequestManageObjectID struct {
	ObjectID primitive.ObjectID `form:"_id" validate:"required"`
}

type RequestManageDelete struct {
	Delete bool `form:"delete"`
}

type RequestManageRoleEnable struct {
	RequestManageObjectID
	RequestManageEnable
}

type RequestManageRoleDelete struct {
	RequestManageObjectID
	RequestManageDelete
}

type RequestManageTokenDelete struct {
	RequestManageObjectID
	Token string `form:"token" validate:"required"`
}

type RequestManageUserCreate struct {
	UserName string `form:"username" validate:"required"`
	RealName string `form:"real_name" validate:"required"`
	Email    string `form:"email" validate:"email"`
	RequestManageEnable
	RequestManagePassword
	Description string `form:"description" validate:"required"`
}

type RequestManageUserID struct {
	UserID primitive.ObjectID `form:"user_id" validate:"required"`
}

type RequestManageUserUpdate struct {
	RequestManageUserID
	UserName string `form:"username" validate:"required"`
	RealName string `form:"real_name" validate:"required"`
	Email    string `form:"email" validate:"email"`
	RequestManageEnable
	Description string `form:"description" validate:"required"`
}

type RequestManageUserEnable struct {
	RequestManageUserID
	RequestManageEnable
}

type RequestManageUserDelete struct {
	RequestManageUserID
	RequestManageDelete
}
