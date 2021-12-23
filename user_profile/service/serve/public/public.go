package public

import (
	"net/http"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/plugins/user_profile/database/model/profile"
	"github.com/haiyiyun/plugins/user_profile/predefined"
	"github.com/haiyiyun/utils/help"
	"github.com/haiyiyun/utils/http/pagination"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/utils/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (self *Service) Route_GET_NicknameAvatar(rw http.ResponseWriter, r *http.Request) {
	var requestUID predefined.RequestServeUserID
	if err := validator.FormStruct(&requestUID, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	profileModel := profile.NewModel(self.M)
	if pf, err := profileModel.GetNickNameAndAvatar(requestUID.UserID); err == nil {
		response.JSON(rw, 0, pf, "")
		return
	}

	response.JSON(rw, http.StatusNotFound, nil, "")
}

func (self *Service) Route_GET_NicknameList(rw http.ResponseWriter, r *http.Request) {
	var requestSN predefined.RequestServeNickname
	if err := validator.FormStruct(&requestSN, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	profileModel := profile.NewModel(self.M)
	filter := profileModel.FilterByNormalProfile()
	filter = append(filter, profileModel.FilterByNicknameWithRegex(requestSN.Nickname)...)

	cnt, _ := profileModel.CountDocuments(r.Context(), filter)
	pg := pagination.Parse(r, cnt)

	projection := bson.D{
		{"_id", 1},
		{"info.nickname", 1},
	}

	opt := options.Find().SetSort(bson.D{
		{"create_time", -1},
	}).SetProjection(projection).SetSkip(pg.SkipNum).SetLimit(pg.PageSize)

	if cur, err := profileModel.Find(r.Context(), filter, opt); err != nil {
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

func (self *Service) Route_GET_SearchNickname(rw http.ResponseWriter, r *http.Request) {
	var requestSSN predefined.RequestServeSearchNickname
	if err := validator.FormStruct(&requestSSN, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	limit := 20

	if requestSSN.Limit > 0 && requestSSN.Limit <= 30 {
		limit = requestSSN.Limit
	}

	profileModel := profile.NewModel(self.M)
	filter := profileModel.FilterByNormalProfile()
	filter = append(filter, profileModel.FilterByNicknameStartWithRegex(requestSSN.Nickname)...)

	projection := bson.D{
		{"_id", 1},
		{"info.nickname", 1},
	}

	opt := options.Find().SetSort(bson.D{
		{"create_time", -1},
	}).SetProjection(projection).SetLimit(int64(limit))

	if cur, err := profileModel.Find(r.Context(), filter, opt); err != nil {
		response.JSON(rw, http.StatusNotFound, nil, "")
	} else {
		items := []help.M{}
		if err := cur.All(r.Context(), &items); err != nil {
			log.Error(err)
			response.JSON(rw, http.StatusServiceUnavailable, nil, "")
		} else {
			response.JSON(rw, 0, items, "")
		}
	}
}
