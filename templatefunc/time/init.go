package templatefunc

import (
	"github.com/haiyiyun/template"
	"github.com/haiyiyun/utils/help"
)

func init() {
	template.AddFunc("UnixToStr", help.NewTime().UnixToStr)
	template.AddFunc("StrToUnix", strToUnix)
}
