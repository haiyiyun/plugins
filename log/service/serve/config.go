package serve

import (
	"github.com/haiyiyun/plugins/log/service/base"
)

type LogConfig struct {
	Log bool `json:"log"`
}

type Config struct {
	base.Config
	LogConfig
}
