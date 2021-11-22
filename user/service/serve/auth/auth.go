package auth

import (
	"net/http"

	"github.com/haiyiyun/plugins/user/database/model"
	"github.com/haiyiyun/plugins/user/database/model/profile"
	"github.com/haiyiyun/plugins/user/database/model/user"
	"github.com/haiyiyun/plugins/user/predefined"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/mongodb/geometry"
	"github.com/haiyiyun/utils/help"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/utils/realip"
	"github.com/haiyiyun/validator"
	"github.com/haiyiyun/validator/form"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (self *Service) Route_POST_Login(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	self.Logout(r)

	var requestLogin predefined.RequestServeAuthLogin

	decoder := form.NewDecoder()
	err := decoder.Decode(&requestLogin, r.Form)
	if err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(requestLogin)
	if err != nil {
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

func (self *Service) Route_POST_TokensByUsernameAndPassword(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var requestUP predefined.RequestServeAuthUsernamePassword

	decoder := form.NewDecoder()
	err := decoder.Decode(&requestUP, r.Form)
	if err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(requestUP)
	if err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	if ts, err := self.GetTokensByUsernameAndPassword(requestUP.Username, requestUP.Password); err == nil {
		response.JSON(rw, 0, ts, "")
	} else {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
	}
}

func (self *Service) Route_DELETE_TokenByUsernameAndPassword(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var requestTUP predefined.RequestServeAuthTokenByUsernameAndPassword

	decoder := form.NewDecoder()
	err := decoder.Decode(&requestTUP, r.Form)
	if err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(requestTUP)
	if err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	tokenID, _ := primitive.ObjectIDFromHex(requestTUP.TokenID)

	if err := self.DeleteTokenByUsernameAndPassword(tokenID, requestTUP.Username, requestTUP.Password); err == nil {
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
	r.ParseForm()

	var requestT predefined.RequestServeAuthTokenID

	decoder := form.NewDecoder()
	err := decoder.Decode(&requestT, r.Form)
	if err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(requestT)
	if err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	tokenID, _ := primitive.ObjectIDFromHex(requestT.TokenID)

	if err := self.DeleteTokenByToken(tokenID, r); err == nil {
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

	decoder := form.NewDecoder()
	err := decoder.Decode(&requestRefresh, r.Form)
	if err != nil {
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

	var requestCheck predefined.RequestServeAuthCheck

	decoder := form.NewDecoder()
	err := decoder.Decode(&requestCheck, r.Form)
	if err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(requestCheck)
	if err != nil {
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

	decoder := form.NewDecoder()
	err := decoder.Decode(&requestCreate, r.Form)
	if err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(requestCreate)
	if err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	var userID primitive.ObjectID
	userID, err = self.CreateUser(r.Context(), requestCreate.Username, requestCreate.Password, requestCreate.Longitude, requestCreate.Latitude, self.Config.EnableProfile, 0)

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

	var requestChangePassword predefined.RequestServeAuthChangePassword

	decoder := form.NewDecoder()
	err := decoder.Decode(&requestChangePassword, r.Form)
	if err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(requestChangePassword)
	if err != nil {
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

	decoder := form.NewDecoder()
	err := decoder.Decode(&requestGuest, r.Form)
	if err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	coordinates := geometry.PointCoordinates{
		requestGuest.Longitude, requestGuest.Latitude,
	}

	userID := primitive.NewObjectID()
	username := userID.Hex()
	usernameMd5 := help.NewString(username).Md5()
	password := help.NewString(usernameMd5).Md5()

	userModel := user.NewModel(self.M)
	ctx := r.Context()
	err = userModel.UseSession(ctx, func(sctx mongo.SessionContext) error {
		if err := sctx.StartTransaction(); err != nil {
			return err
		}

		u := model.User{
			ID:       userID,
			Name:     username,
			Password: password,
			Guest:    true,
			Enable:   true,
		}

		if coordinates != geometry.NilPointCoordinates {
			u.Location = geometry.NewPoint(coordinates)
		}

		_, err := userModel.Create(r.Context(), u)
		if err != nil {
			sctx.AbortTransaction(sctx)
			return err
		}

		profileModel := profile.NewModel(self.M)
		_, err = profileModel.Create(r.Context(), model.Profile{
			UserID: userID,
			Enable: true,
		})

		if err != nil {
			sctx.AbortTransaction(sctx)
			return err
		}

		return sctx.CommitTransaction(sctx)
	})

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

	decoder := form.NewDecoder()
	err := decoder.Decode(&requestGuestToUser, r.Form)
	if err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	validate := validator.New()
	err = validate.Struct(requestGuestToUser)
	if err != nil {
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
