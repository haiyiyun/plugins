package predefined

type RequestServeNickname struct {
	Nickname string `form:"nickname" validate:"required"`
}

type RequestServeAvatar struct {
	Avatar string `form:"avatar" validate:"required"`
}

type RequestServePhotos struct {
	Photos []string `form:"photos" validate:"gt=0,dive,required"`
}

type RequestServeTags struct {
	Tags []string `form:"tags" validate:"gt=0,dive,required"`
}

type RequestServeCoverImage struct {
	Image string `form:"image" validate:"required"`
}

type RequestServeCoverVideo struct {
	Video string `form:"video" validate:"required"`
}

type RequestServeCoverVoice struct {
	Voice string `form:"voice" validate:"required"`
}
