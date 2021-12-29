package predefined

import (
	"github.com/haiyiyun/utils/help"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RequestServeAuthUsername struct {
	Username string `form:"username" validate:"required"`
}

type RequestServeAuthPassword struct {
	Password string `form:"password" validate:"required"`
}

type RequestServeAuthSecurePassword struct {
	SecurePassword string `form:"secure_password" validate:"required"`
}

type RequestServeAuthUsernamePassword struct {
	RequestServeAuthUsername
	RequestServeAuthPassword
}

type RequestServeLongitudeLatitude struct {
	Longitude float64 `form:"longitude,omitempty" validate:"numeric"` //经度
	Latitude  float64 `form:"latitude,omitempty" validate:"numeric"`  //维度
}

type RequestServeDistance struct {
	MaxDistance float64 `form:"max_distance,omitempty" validate:"numeric"`
	MinDistance float64 `form:"min_distance,omitempty" validate:"numeric"`
}

type RequestServeOnlineLongitudeLatitude struct {
	OnlineLongitude float64 `form:"pnline_longitude,omitempty" validate:"numeric"` //经度
	OnlineLatitude  float64 `form:"pnline_latitude,omitempty" validate:"numeric"`  //维度
}

type RequestServeOnlineDistance struct {
	OnlineMaxDistance float64 `form:"pnline_max_distance,omitempty" validate:"numeric"`
	OnlineMinDistance float64 `form:"pnline_min_distance,omitempty" validate:"numeric"`
}

type RequestServeAuthLogin struct {
	RequestServeAuthUsernamePassword
	RequestServeLongitudeLatitude
}

type RequestServeAuthRefresh struct {
	RequestServeLongitudeLatitude
}

type RequestServeAuthCreate struct {
	RequestServeAuthUsernamePassword
	RequestServeLongitudeLatitude
}

type RequestServeAuthGuest struct {
	RequestServeLongitudeLatitude
}

type RequestServeAuthGuestToUser struct {
	RequestServeAuthUsernamePassword
}

type RequestServeAuthTokenID struct {
	TokenID primitive.ObjectID `form:"token_id" validate:"required"`
}

type RequestServeAuthTokenByUsernameAndPassword struct {
	RequestServeAuthTokenID
	RequestServeAuthUsernamePassword
}

type RequestServeUserList struct {
	UserID        primitive.ObjectID `form:"user_id"`
	ExtensionID   int                `form:"extension_id"`
	GuestQuery    bool               `form:"guest_query"`
	Guest         bool               `form:"guest"`
	Roles         []string           `form:"roles"`
	RolesWithTime help.DateTime      `form:"roles_with_time"`
	Tags          []string           `form:"tags"`
	TagsWithTime  help.DateTime      `form:"tags_with_time"`
	Level         int                `form:"level"`
	GteLevel      int                `form:"gte_level"`
	LteLevel      int                `form:"lte_level"`
	OnlineQuery   bool               `form:"online_query"`
	Online        bool               `form:"online"`
	RequestServeLongitudeLatitude
	RequestServeDistance
	RequestServeOnlineLongitudeLatitude
	RequestServeOnlineDistance
}
