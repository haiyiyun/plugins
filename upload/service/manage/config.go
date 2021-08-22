package manage

import (
	"github.com/haiyiyun/plugins/upload/service/base"
)

type UploadConfig struct {
	WebRouter                 bool   `json:"web_router"`
	WebRouterRootPath         string `json:"web_router_root_path"`
	BuildInFileServer         bool   `json:"build_in_file_server"`
	DownloadLocalUrlDirectory string `json:"download_local_url_directory"`
	AllowDownloadLocal        bool   `json:"allow_download_local"`
	CheckUser                 bool   `json:"check_user"`
}

type Config struct {
	base.Config
	UploadConfig
}
