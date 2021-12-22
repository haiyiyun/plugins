package user

import (
	"net/http"
	"strconv"
	"time"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/plugins/user/database/model/user"
	"github.com/haiyiyun/plugins/user/predefined"

	"github.com/haiyiyun/mongodb/geometry"
	"github.com/haiyiyun/utils/help"
	"github.com/haiyiyun/utils/http/pagination"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/utils/realip"
	"github.com/haiyiyun/utils/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// valid := validator.Validation{}
	// valid.Float(longitude).Key("longitude").Message("longitude必须是正确的经度坐标点")
	// valid.Float(latitude).Key("latitude").Message("latitude必须是正确的纬度坐标点")

	// if valid.HasErrors() {
	// 	response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
	// 	return
	// }

	ip := realip.RealIP(r)

	coordinates := geometry.PointCoordinates{
		longitude, latitude,
	}

	change := bson.D{
		{"online.online", true},
		{"online.ip", ip},
		{"online.location", geometry.NewPoint(coordinates)},
		{"online.online_time", time.Now()},
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

	// valid := validator.Validation{}
	// valid.Float(longitude).Key("longitude").Message("longitude必须是正确的经度坐标点")
	// valid.Float(latitude).Key("latitude").Message("latitude必须是正确的纬度坐标点")

	// if valid.HasErrors() {
	// 	response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
	// 	return
	// }

	ip := realip.RealIP(r)

	coordinates := geometry.PointCoordinates{
		longitude, latitude,
	}

	change := bson.D{
		{"online.online", false},
		{"online.ip", ip},
		{"online.location", geometry.NewPoint(coordinates)},
		{"online.offline_time", time.Now()},
	}

	userModel := user.NewModel(self.M)
	if ur, err := userModel.Set(r.Context(), userModel.FilterByID(u.ID), change); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_GET_List(rw http.ResponseWriter, r *http.Request) {
	_, found := self.GetUserInfo(r)
	if !found {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var requestSUL predefined.RequestServeUserList
	if err := validator.FormStruct(&requestSUL, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	userModel := user.NewModel(self.M)
	filter := userModel.FilterByNormalUser()

	if requestSUL.GuestQuery {
		filter = append(filter, userModel.FilterByGuest(requestSUL.Guest)...)
	}

	if requestSUL.OnlineQuery {
		filter = append(filter, userModel.FilterByOnline(requestSUL.Online)...)
	}

	if requestSUL.Level > 0 {
		filter = append(filter, userModel.FilterByLevel(requestSUL.Level)...)
	}

	if requestSUL.GteLevel > 0 {
		filter = append(filter, userModel.FilterByGteLevel(requestSUL.GteLevel)...)
	}

	if requestSUL.LteLevel > 0 {
		filter = append(filter, userModel.FilterByLteLevel(requestSUL.LteLevel)...)
	}

	if len(requestSUL.Roles) > 0 {
		if !requestSUL.RolesWithTime.Time.IsZero() {
			filter = append(filter, userModel.FilterByRolesWithTime(requestSUL.Roles, requestSUL.RolesWithTime.Time)...)
		} else {
			filter = append(filter, userModel.FilterByRoles(requestSUL.Roles)...)
		}
	}

	if len(requestSUL.Tags) > 0 {
		if !requestSUL.TagsWithTime.Time.IsZero() {
			filter = append(filter, userModel.FilterByTagsWithTime(requestSUL.Tags, requestSUL.TagsWithTime.Time)...)
		} else {
			filter = append(filter, userModel.FilterByTags(requestSUL.Tags)...)
		}
	}

	coordinates := geometry.PointCoordinates{
		requestSUL.Longitude, requestSUL.Latitude,
	}

	if coordinates != geometry.NilPointCoordinates {
		filter = append(filter, userModel.FilterByLocation(geometry.NewPoint(coordinates), requestSUL.MaxDistance, requestSUL.MinDistance)...)
	}

	onlineCoordinates := geometry.PointCoordinates{
		requestSUL.OnlineLongitude, requestSUL.OnlineLatitude,
	}

	if onlineCoordinates != geometry.NilPointCoordinates {
		filter = append(filter, userModel.FilterByOnlineLocation(geometry.NewPoint(onlineCoordinates), requestSUL.OnlineMaxDistance, requestSUL.OnlineMinDistance)...)
	}

	cnt, _ := userModel.CountDocuments(r.Context(), filter)
	pg := pagination.Parse(r, cnt)

	projection := bson.D{
		{"_id", 1},
		{"extension_id", 1},
	}

	opt := options.Find().SetSort(bson.D{
		{"create_time", -1},
	}).SetProjection(projection).SetSkip(pg.SkipNum).SetLimit(pg.PageSize)

	if cur, err := userModel.Find(r.Context(), filter, opt); err != nil {
		response.JSON(rw, http.StatusNotFound, nil, "")
	} else {
		items := []help.M{}
		if err := cur.All(r.Context(), &items); err != nil {
			log.Error(err)
			response.JSON(rw, http.StatusServiceUnavailable, nil, "")
		} else {
			rpr := response.ResponsePaginationResult{
				Total: cnt,
				Items: items,
			}

			response.JSON(rw, 0, rpr, "")
		}
	}

}
