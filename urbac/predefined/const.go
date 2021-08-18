package predefined

const (
	TokenTypeSelf  = "self"
	TokenTypeAuth  = "auth"
	TokenTypeGuise = "guise"
)

const (
	RoleRightTypeAction = iota
	RoleRightTypeModule
	RoleRightTypeApp
	RoleRightTypePlatform
)

const (
	ApplicationTypeCode    = "code"
	ApplicationTypeVirtual = "virtual"
)

const (
	ApplicationLevelApp    = "app"
	ApplicationLevelModule = "module"
	ApplicationLevelAction = "action"
	ApplicationLevelMethod = "method"
)

const (
	StatusCodeLoginLimit     = 5031
	StatusCodeLoginLimitText = "登录限制"
)
