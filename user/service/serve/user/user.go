package user

import (
	"net/http"
	"strconv"
	"time"

	"github.com/haiyiyun/plugins/user/database/model/user"

	"github.com/haiyiyun/mongodb/geometry"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/utils/realip"
	"github.com/haiyiyun/validator"
	"go.mongodb.org/mongo-driver/bson"
)

func (self *Service) Route_POST_Online(rw http.ResponseWriter, r *http.Request) {
	u, found := self.GetUserInfo(r)
	if !found {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()
	longitudeStr := r.FormValue("longitude") //经度
	latitudeStr := r.FormValue("latitude")   //维度
	longitude, _ := strconv.ParseFloat(longitudeStr, 64)
	latitude, _ := strconv.ParseFloat(latitudeStr, 64)

	valid := validator.Validation{}
	valid.Float(longitude).Key("longitude").Message("longitude必须是正确的经度坐标点")
	valid.Float(latitude).Key("latitude").Message("latitude必须是正确的纬度坐标点")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	ip := realip.RealIP(r)

	coordinates := geometry.PointCoordinates{
		longitude, latitude,
	}

	change := bson.D{
		{"online", bson.D{
			{"online", true},
			{"ip", ip},
			{"location", geometry.NewPoint(coordinates)},
			{"online_time", time.Now()},
		}},
	}

	userModel := user.NewModel(self.M)
	if ur, err := userModel.Set(r.Context(), userModel.FilterByID(u.ID), change); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_POST_Offline(rw http.ResponseWriter, r *http.Request) {
	u, found := self.GetUserInfo(r)
	if !found {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()
	longitudeStr := r.FormValue("longitude") //经度
	latitudeStr := r.FormValue("latitude")   //维度
	longitude, _ := strconv.ParseFloat(longitudeStr, 64)
	latitude, _ := strconv.ParseFloat(latitudeStr, 64)

	valid := validator.Validation{}
	valid.Float(longitude).Key("longitude").Message("longitude必须是正确的经度坐标点")
	valid.Float(latitude).Key("latitude").Message("latitude必须是正确的纬度坐标点")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	ip := realip.RealIP(r)

	coordinates := geometry.PointCoordinates{
		longitude, latitude,
	}

	change := bson.D{
		{"online", bson.D{
			{"online", false},
			{"ip", ip},
			{"location", geometry.NewPoint(coordinates)},
			{"offline_time", time.Now()},
		}},
	}

	userModel := user.NewModel(self.M)
	if ur, err := userModel.Set(r.Context(), userModel.FilterByID(u.ID), change); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}
