package templatefunc

import (
	"time"

	"github.com/haiyiyun/utils/help"
)

func strToUnix(s string) int64 {
	t := help.NewTime()

	if s == "now" {
		return time.Now().Unix()
	}

	return t.StrToUnix(s)
}
