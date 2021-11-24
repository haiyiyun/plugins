package serve

import (
	"github.com/haiyiyun/plugins/upload/service/base"
)

type UploadConfig struct {
	WebRouter                 bool   `json:"web_router"`
	WebRouterRootPath         string `json:"web_router_root_path"`
	BuildInFileServer         bool   `json:"build_in_file_server"`
	PublishDownloadUrl        bool   `json:"publish_download_url"`
	DownloadLocalUrlDirectory string `json:"download_local_url_directory"`
	AllowDownloadLocal        bool   `json:"allow_download_local"`
	MaxUploadFileSize         int64  `json:"max_file_size"`
	CheckUser                 bool   `json:"check_user"`
}

type Config struct {
	base.Config
	UploadConfig
}
