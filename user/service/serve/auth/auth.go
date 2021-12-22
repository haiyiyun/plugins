package auth

import (
	"net/http"

	"github.com/haiyiyun/plugins/user/database/model/user"
	"github.com/haiyiyun/plugins/user/predefined"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/mongodb/geometry"
	"github.com/haiyiyun/utils/help"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/utils/realip"
	"github.com/haiyiyun/utils/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (self *Service) Route_POST_Login(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	self.Logout(r)

	var requestLogin predefined.RequestServeAuthLogin
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

func (self *Service) Route_POST_GetTokens(rw http.ResponseWriter, r *http.Request) {
	if self.Config.ProhibitDeleteTokensByUsernamePassword {
		response.JSON(rw, http.StatusForbidden, nil, "")
		return
	}

	r.ParseForm()

	var requestUP predefined.RequestServeAuthUsernamePassword
	if err := validator.FormStruct(&requestUP, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	if ts, err := self.GetTokensByUsernameAndPassword(requestUP.Username, requestUP.Password); err == nil {
		response.JSON(rw, 0, ts, "")
	} else {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
	}
}

func (self *Service) Route_POST_DeleteToken(rw http.ResponseWriter, r *http.Request) {
	if self.Config.ProhibitDeleteTokensByUsernamePassword {
		response.JSON(rw, http.StatusForbidden, nil, "")
		return
	}

	r.ParseForm()

	var requestTUP predefined.RequestServeAuthTokenByUsernameAndPassword
	if err := validator.FormStruct(&requestTUP, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := self.DeleteTokenByUsernameAndPassword(requestTUP.TokenID, requestTUP.Username, requestTUP.Password); err == nil {
		response.JSON(rw, 0, nil, "")
	} else {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
	}
}

func (self *Service) Route_GET_Tokens(rw http.ResponseWriter, r *http.Request) {
	if ts, err := self.GetTokensByToken(r); err == nil {
		response.JSON(rw, 0, ts, "")
	} else {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
	}
}

func (self *Service) Route_DELETE_Token(rw http.ResponseWriter, r *http.Request) {
	var requestT predefined.RequestServeAuthTokenID
	if err := validator.FormStruct(&requestT, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := self.DeleteTokenByToken(requestT.TokenID, r); err == nil {
		response.JSON(rw, 0, nil, "")
	} else {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
	}
}

func (self *Service) Route_POST_Refresh(rw http.ResponseWriter, r *http.Request) {
	u, found := self.GetUserInfo(r)
	if !found {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	self.Logout(r)

	var requestRefresh predefined.RequestServeAuthRefresh
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
			log.Debug(err)
			response.JSON(rw, http.StatusUnauthorized, nil, "")
		}
	} else {
		response.JSON(rw, 0, m, "")
	}
}

func (self *Service) Route_GET_Logout(rw http.ResponseWriter, r *http.Request) {
	_, found := self.GetUserInfo(r)
	if !found {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	self.Logout(r)

	response.JSON(rw, 0, nil, "")

}

func (self *Service) Route_GET_GetUserInfo(rw http.ResponseWriter, r *http.Request) {
	if u, found := self.GetUserInfo(r); found {
		response.JSON(rw, 0, u, "")
	} else {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
	}
}

func (self *Service) Route_GET_Check(rw http.ResponseWriter, r *http.Request) {
	if self.Config.ProhibitCreateUser {
		response.JSON(rw, http.StatusForbidden, nil, "")
		return
	}

	r.ParseForm()

	var requestCheck predefined.RequestServeAuthUsername
	if err := validator.FormStruct(&requestCheck, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	if cnt, err := self.CheckUser(r.Context(), requestCheck.Username); err == nil {
		if cnt == 0 {
			response.JSON(rw, 0, help.M{
				"exist": false,
			}, "不存在")
		} else {
			response.JSON(rw, 0, help.M{
				"exist": true,
			}, "存在")
		}
	} else {
		log.Error(err)
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	}
}

func (self *Service) Route_POST_Create(rw http.ResponseWriter, r *http.Request) {
	if self.Config.ProhibitCreateUser {
		response.JSON(rw, http.StatusForbidden, nil, "")
		return
	}

	r.ParseForm()

	var requestCreate predefined.RequestServeAuthCreate
	if err := validator.FormStruct(&requestCreate, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	var userID primitive.ObjectID
	userID, err := self.CreateUser(r.Context(), primitive.NilObjectID, requestCreate.Username, requestCreate.Password, requestCreate.Longitude, requestCreate.Latitude, self.Config.EnableProfile, 0, false)

	if err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, help.M{
			"user_id": userID,
		}, "")
	}
}

func (self *Service) Route_POST_ChangePassword(rw http.ResponseWriter, r *http.Request) {
	u, found := self.GetUserInfo(r)
	if !found {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestChangePassword predefined.RequestServeAuthPassword
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

//username => user_id
//password => md5(username)
func (self *Service) Route_POST_Guest(rw http.ResponseWriter, r *http.Request) {
	if self.Config.ProhibitCreateUser {
		response.JSON(rw, http.StatusForbidden, nil, "")
		return
	}

	r.ParseForm()

	var requestGuest predefined.RequestServeAuthGuest
	if err := validator.FormStruct(&requestGuest, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	userID := primitive.NewObjectID()
	username := userID.Hex()
	usernameMd5 := help.NewString(username).Md5()
	password := help.NewString(usernameMd5).Md5()

	userID, err := self.CreateUser(r.Context(), userID, username, password, requestGuest.Longitude, requestGuest.Latitude, self.Config.EnableProfile, 0, false)

	if err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, help.M{
			"user_id": userID, //调用guset方需要将user_id永久保存,作为guest的username登录
		}, "")
	}
}

func (self *Service) Route_POST_GuestToUser(rw http.ResponseWriter, r *http.Request) {
	if self.Config.ProhibitCreateUser {
		response.JSON(rw, http.StatusForbidden, nil, "")
		return
	}

	u, found := self.GetUserInfo(r)
	if !found {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestGuestToUser predefined.RequestServeAuthGuestToUser
	if err := validator.FormStruct(&requestGuestToUser, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	userModel := user.NewModel(self.M)
	if err := userModel.GuestToUser(u.ID, requestGuestToUser.Username, requestGuestToUser.Password); err == nil {
		response.JSON(rw, 0, nil, "")
	} else {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	}
}
