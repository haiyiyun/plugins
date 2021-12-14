package predefined

type RequestServeLookup struct {
	Key       string `form:"key" validate:"required"`
	Structure string `form:"structure" validate:"required,oneof=slice,map"`
	Want      string `form:"want" validate:"required,oneof=lable,value,lablevalue"`
	Flush     bool   `form:"flush"`
}
