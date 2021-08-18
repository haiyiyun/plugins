package templatefunc

import (
	"github.com/haiyiyun/template"
)

func init() {
	template.AddFunc("OptionsForSelect", optionsForSelect)
}
