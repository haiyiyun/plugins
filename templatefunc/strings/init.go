package strings

import (
	"strconv"
	"strings"

	"github.com/haiyiyun/template"
)

func init() {
	template.AddFunc("SubStr", substr)
	template.AddFunc("HasSuffix", strings.HasSuffix)
	template.AddFunc("IntToString", func(i int) string {
		return strconv.FormatInt(int64(i), 10)
	})
}
