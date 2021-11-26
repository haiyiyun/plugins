package cors

type CorsConfig struct {
	AccessControlAllowOrigin      string `json:"access_control_allow_origin"`
	AccessControlAllowHeaders     string `json:"access_control_allow_headers"`
	AccessControlAllowMethods     string `json:"access_control_allow_methods"`
	AccessControlExposeHeaders    string `json:"access_control_expose_headers"`
	AccessControlAllowCredentials string `json:"access_control_allow_credentials"`
}

type Config struct {
	CorsConfig
}
