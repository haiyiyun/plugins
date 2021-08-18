package operator

import "github.com/haiyiyun/template"

func init() {
	template.AddFunc("op", operator)
}
