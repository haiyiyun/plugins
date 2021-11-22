package predefined

type RequestServeAuthUsernamePassword struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type RequestServeAuthLongitudeLatitude struct {
	Longitude float64 `form:",omitempty"` //经度
	Latitude  float64 `form:",omitempty"` //维度
}

type RequestServeAuthLogin struct {
	RequestServeAuthUsernamePassword
	RequestServeAuthLongitudeLatitude
}

type RequestServeAuthRefresh struct {
	RequestServeAuthLongitudeLatitude
}

type RequestServeAuthCheck struct {
	Username string `validate:"required"`
}

type RequestServeAuthCreate struct {
	RequestServeAuthUsernamePassword
	RequestServeAuthLongitudeLatitude
}

type RequestServeAuthChangePassword struct {
	Password string `validate:"required"`
}

type RequestServeAuthGuest struct {
	RequestServeAuthLongitudeLatitude
}

type RequestServeAuthGuestToUser struct {
	RequestServeAuthUsernamePassword
}

type RequestServeAuthTokenID struct {
	TokenID string `validate:"required,bson_object_id"`
}

type RequestServeAuthTokenByUsernameAndPassword struct {
	RequestServeAuthTokenID
	RequestServeAuthUsernamePassword
}
