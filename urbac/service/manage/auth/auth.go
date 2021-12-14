package auth

import (
	"net/http"

	"github.com/haiyiyun/mongodb/geometry"
	"github.com/haiyiyun/plugins/urbac/database/model/user"
	"github.com/haiyiyun/plugins/urbac/predefined"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/utils/realip"
	"github.com/haiyiyun/utils/validator"
)

func (self *Service) Route_POST_Login(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	self.Logout(r)

	var requestLogin predefined.RequestManageLogin
	if err := validator.FormStruct(&requestLogin, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	coordinates := geometry.PointCoordinates{
		requestLogin.Longitude, requestLogin.Latitude,
	}

	if m, err := self.Login(requestLogin.Username, requestLogin.Password, realip.RealIP(r), r.Header.Get("User-Agent"), coordinates); err != nil {
		if err.Error() == predefined.StatusCodeLoginLimitText {
			response.JSON(rw, predefined.StatusCodeLoginLimit, nil, predefined.StatusCodeLoginLimitText)
		} else {
			response.JSON(rw, http.StatusUnauthorized, nil, "")
		}
	} else {
		response.JSON(rw, 0, m, "")
	}
}

func (self *Service) Route_POST_Refresh(rw http.ResponseWriter, r *http.Request) {
	if u, found := self.GetUserInfo(r); found {
		r.ParseForm()
		self.Logout(r)

		var requestRefresh predefined.RequestManageRefresh
		if err := validator.FormStruct(&requestRefresh, r.Form); err != nil {
			response.JSON(rw, http.StatusBadRequest, nil, err.Error())
			return
		}

		coordinates := geometry.PointCoordinates{
			requestRefresh.Longitude, requestRefresh.Latitude,
		}

		if m, err := self.CreateToken(r.Context(), u, realip.RealIP(r), r.Header.Get("User-Agent"), coordinates); err != nil {
			if err.Error() == predefined.StatusCodeLoginLimitText {
				response.JSON(rw, predefined.StatusCodeLoginLimit, nil, predefined.StatusCodeLoginLimitText)
			} else {
				response.JSON(rw, http.StatusUnauthorized, nil, "")
			}
		} else {
			response.JSON(rw, 0, m, "")
		}
	} else {
		response.JSON(rw, http.StatusNotFound, nil, "")
	}
}

func (self *Service) Route_GET_Logout(rw http.ResponseWriter, r *http.Request) {
	self.Logout(r)

	response.JSON(rw, 0, nil, "")

}

func (self *Service) Route_POST_ChangePassword(rw http.ResponseWriter, r *http.Request) {
	u, found := self.GetUserInfo(r)
	if !found {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestChangePassword predefined.RequestManagePassword
	if err := validator.FormStruct(&requestChangePassword, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	userModel := user.NewModel(self.M)
	if err := userModel.ChangePassword(u.ID, requestChangePassword.Password); err == nil {
		response.JSON(rw, 0, nil, "")
	} else {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	}

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
