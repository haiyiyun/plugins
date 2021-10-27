package auth

import (
	"net/http"
	"strconv"

	"github.com/haiyiyun/plugins/urbac/database/model/user"
	"github.com/haiyiyun/plugins/urbac/predefined"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/utils/realip"
	"github.com/haiyiyun/validator"
)

func (self *Service) Route_POST_Login(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	self.Logout(r)
	username := r.FormValue("username")
	password := r.FormValue("password")
	geoLongitudeStr := r.FormValue("longitude") //经度
	geoLatitudeStr := r.FormValue("latitude")   //维度
	geoLongitude, _ := strconv.ParseFloat(geoLongitudeStr, 64)
	geoLatitude, _ := strconv.ParseFloat(geoLatitudeStr, 64)
	geo := [2]float64{
		geoLongitude, geoLatitude,
	}

	if m, err := self.Login(username, password, realip.RealIP(r), r.Header.Get("User-Agent"), geo); err != nil {
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
		geoLongitudeStr := r.FormValue("longitude") //经度
		geoLatitudeStr := r.FormValue("latitude")   //维度
		geoLongitude, _ := strconv.ParseFloat(geoLongitudeStr, 64)
		geoLatitude, _ := strconv.ParseFloat(geoLatitudeStr, 64)
		geo := [2]float64{
			geoLongitude, geoLatitude,
		}

		if m, err := self.CreateToken(r.Context(), u, realip.RealIP(r), r.Header.Get("User-Agent"), geo); err != nil {
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
	password := r.FormValue("password")

	valid := validator.Validation{}
	valid.Required(password).Key("password").Message("password不能为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	userModel := user.NewModel(self.M)
	if err := userModel.ChangePassword(u.ID, password); err == nil {
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
