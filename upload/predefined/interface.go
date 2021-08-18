package predefined

import (
	"net/http"

	"github.com/haiyiyun/plugins/upload/database/model"
)

type Upload interface {
	SaveEncodeFile(*http.Request, string) (*model.Upload, error)
	SaveFormFile(*http.Request, string) (*model.Upload, error)
}
