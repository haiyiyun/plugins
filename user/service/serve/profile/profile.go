package profile

import (
	"net/http"

	"github.com/haiyiyun/plugins/user/database/model/profile"

	"github.com/haiyiyun/utils/http/response"
)

func (self *Service) Route_GET_Profile(rw http.ResponseWriter, r *http.Request) {
	if u, found := self.GetUserInfo(r); found {
		profileModel := profile.NewModel(self.M)

		if pf, err := profileModel.GetInfo(u.ID); err == nil {
			response.JSON(rw, 0, pf, "")
			return
		}
	}

	response.JSON(rw, http.StatusNotFound, nil, "")
}
