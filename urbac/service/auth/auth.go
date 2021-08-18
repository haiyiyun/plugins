package auth

import (
	"net/http"

	"github.com/haiyiyun/plugins/urbac/predefined"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/utils/realip"
)

func (self *Service) Route_POST_Login(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	self.Logout(r)
	username := r.FormValue("username")
	password := r.FormValue("password")
	if m, err := self.Login(username, password, realip.RealIP(r), r.Header.Get("User-Agent")); err != nil {
		if err.Error() == predefined.StatusCodeLoginLimitText {
			response.JSON(rw, predefined.StatusCodeLoginLimit, nil, predefined.StatusCodeLoginLimitText)
		} else {
			response.JSON(rw, http.StatusUnauthorized, nil, "")
		}
	} else {
		response.JSON(rw, 0, m, "")
	}
}

func (self *Service) Route_GET_Logout(rw http.ResponseWriter, r *http.Request) {
	self.Logout(r)

	response.JSON(rw, 0, nil, "")

}

func (self *Service) Route_GET_GetUserInfo(rw http.ResponseWriter, r *http.Request) {
	if u, found := self.GetUserInfo(r); found {
		if u.Setting.HomePath == "" && self.Config.DefaultHomePath != "" {
			u.Setting.HomePath = self.Config.DefaultHomePath
		}

		response.JSON(rw, 0, u, "")
	} else {
		response.JSON(rw, http.StatusNotFound, nil, "")
	}
}

func (self *Service) Route_GET_GetRouteList(rw http.ResponseWriter, r *http.Request) {
	if u, found := self.GetUserInfo(r); found {
		apps := self.GetApplications(u.ID, self.Config.CheckRight)

		response.JSON(rw, 0, apps["route"], "")
	} else {
		response.JSON(rw, http.StatusNotFound, nil, "")
	}
}

func (self *Service) Route_GET_GetPermissionCode(rw http.ResponseWriter, r *http.Request) {
	if u, found := self.GetUserInfo(r); found {
		apps := self.GetApplications(u.ID, self.Config.CheckRight)

		response.JSON(rw, 0, apps["permission_code"], "")
	} else {
		response.JSON(rw, http.StatusNotFound, nil, "")
	}
}
