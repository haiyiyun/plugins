package predefined

type RequestServeAuthLogin struct {
	Username  string  `validate:"required"`
	Password  string  `validate:"required"`
	Longitude float64 `form:",omitempty"` //经度
	Latitude  float64 `form:",omitempty"` //维度
}

type RequestServeAuthRefresh struct {
	Longitude float64 `form:",omitempty"` //经度
	Latitude  float64 `form:",omitempty"` //维度
}

type RequestServeAuthCheck struct {
	Username string `validate:"required"`
}

type RequestServeAuthCreate struct {
	Username  string  `validate:"required"`
	Password  string  `validate:"required"`
	Longitude float64 `form:",omitempty"` //经度
	Latitude  float64 `form:",omitempty"` //维度
}

type RequestServeAuthChangePassword struct {
	Password string `validate:"required"`
}

type RequestServeAuthGuest struct {
	Longitude float64 `form:",omitempty"` //经度
	Latitude  float64 `form:",omitempty"` //维度
}

type RequestServeAuthGuestToUser struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}
