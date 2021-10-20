package auth

import (
	"net/http"
	"strconv"

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
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (self *Service) Route_POST_Login(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	self.Logout(r)
	username := r.FormValue("username")
	password := r.FormValue("password")
	longitudeStr := r.FormValue("longitude") //经度
	latitudeStr := r.FormValue("latitude")   //维度
	longitude, _ := strconv.ParseFloat(longitudeStr, 64)
	latitude, _ := strconv.ParseFloat(latitudeStr, 64)
	coordinates := geometry.PointCoordinates{
		longitude, latitude,
	}

	if m, err := self.Login(username, password, realip.RealIP(r), r.Header.Get("User-Agent"), coordinates); err != nil {
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
		longitudeStr := r.FormValue("longitude") //经度
		latitudeStr := r.FormValue("latitude")   //维度
		longitude, _ := strconv.ParseFloat(longitudeStr, 64)
		latitude, _ := strconv.ParseFloat(latitudeStr, 64)
		coordinates := geometry.PointCoordinates{
			longitude, latitude,
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
	} else {
		response.JSON(rw, http.StatusNotFound, nil, "")
	}
}

func (self *Service) Route_GET_Logout(rw http.ResponseWriter, r *http.Request) {
	self.Logout(r)

	response.JSON(rw, 0, nil, "")

}

func (self *Service) Route_GET_GetUserInfo(rw http.ResponseWriter, r *http.Request) {
	if u, found := self.GetUserInfo(r); found {
		response.JSON(rw, 0, u, "")
	} else {
		response.JSON(rw, http.StatusNotFound, nil, "")
	}
}

func (self *Service) Route_GET_Check(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")

	valid := validator.Validation{}
	valid.Required(username).Key("username").Message("username不能为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	userModel := user.NewModel(self.M)
	if cnt, err := userModel.CountDocuments(r.Context(), userModel.FilterByName(username)); err == nil {
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
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	longitudeStr := r.FormValue("longitude") //经度
	latitudeStr := r.FormValue("latitude")   //维度
	longitude, _ := strconv.ParseFloat(longitudeStr, 64)
	latitude, _ := strconv.ParseFloat(latitudeStr, 64)
	coordinates := geometry.PointCoordinates{
		longitude, latitude,
	}

	valid := validator.Validation{}
	valid.Required(username).Key("username").Message("username不能为空")
	valid.Required(password).Key("password").Message("password不能为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	userModel := user.NewModel(self.M)
	ctx := r.Context()
	var userID primitive.ObjectID
	err := userModel.UseSession(ctx, func(sctx mongo.SessionContext) error {
		u := model.User{
			Name:     username,
			Password: help.NewString(password).Md5(),
			Enable:   true,
		}

		if coordinates != geometry.NilPointCoordinates {
			u.Location = geometry.NewPoint(coordinates)
		}

		ior, err := userModel.Create(r.Context(), u)

		if err != nil {
			sctx.AbortTransaction(sctx)
			log.Error("Create user error:", err)
			return err
		}

		userID = ior.InsertedID.(primitive.ObjectID)

		if self.Config.EnableProfile {
			profileModel := profile.NewModel(self.M)
			_, err = profileModel.Create(r.Context(), model.Profile{
				UserID: userID,
				Enable: true,
			})

			if err != nil {
				sctx.AbortTransaction(sctx)
				log.Error("Create profile error:", err)
				return err
			}
		}

		sctx.CommitTransaction(sctx)
		return err
	})

	if err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, help.M{
			"user_id": userID,
		}, "")
	}
}

func (self *Service) Route_POST_Guest(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	longitudeStr := r.FormValue("longitude") //经度
	latitudeStr := r.FormValue("latitude")   //维度
	longitude, _ := strconv.ParseFloat(longitudeStr, 64)
	latitude, _ := strconv.ParseFloat(latitudeStr, 64)
	coordinates := geometry.PointCoordinates{
		longitude, latitude,
	}

	userID := primitive.NewObjectID()
	username := "访客" + userID.Hex()
	password := help.NewString("").Md5()

	userModel := user.NewModel(self.M)
	ctx := r.Context()
	err := userModel.UseSession(ctx, func(sctx mongo.SessionContext) error {
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
			log.Error("Create user error:", err)
			return err
		}

		profileModel := profile.NewModel(self.M)
		_, err = profileModel.Create(r.Context(), model.Profile{
			UserID: userID,
			Enable: true,
		})

		if err != nil {
			sctx.AbortTransaction(sctx)
			log.Error("Create profile error:", err)
			return err
		}

		sctx.CommitTransaction(sctx)
		return err
	})

	if err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, help.M{
			"user_id": userID,
		}, "")
	}
}
