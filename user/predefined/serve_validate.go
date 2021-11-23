package predefined

type RequestServeAuthUsername struct {
	Username string `form:"username" validate:"required"`
}

type RequestServeAuthPassword struct {
	Password string `form:"password" validate:"required"`
}

type RequestServeAuthUsernamePassword struct {
	RequestServeAuthUsername
	RequestServeAuthPassword
}

type RequestServeAuthLongitudeLatitude struct {
	Longitude float64 `form:"longitude,omitempty"` //经度
	Latitude  float64 `form:"latitude,omitempty"`  //维度
}

type RequestServeAuthLogin struct {
	RequestServeAuthUsernamePassword
	RequestServeAuthLongitudeLatitude
}

type RequestServeAuthRefresh struct {
	RequestServeAuthLongitudeLatitude
}

type RequestServeAuthCreate struct {
	RequestServeAuthUsernamePassword
	RequestServeAuthLongitudeLatitude
}

type RequestServeAuthGuest struct {
	RequestServeAuthLongitudeLatitude
}

type RequestServeAuthGuestToUser struct {
	RequestServeAuthUsernamePassword
}

type RequestServeAuthTokenID struct {
	TokenID string `form:"token_id" validate:"required,bson_object_id"`
}

type RequestServeAuthTokenByUsernameAndPassword struct {
	RequestServeAuthTokenID
	RequestServeAuthUsernamePassword
}
