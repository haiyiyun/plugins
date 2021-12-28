package follow

import (
	"net/http"
	"time"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/plugins/content/database/model/follow_content"
	"github.com/haiyiyun/plugins/content/predefined"
	"github.com/haiyiyun/utils/help"
	"github.com/haiyiyun/utils/http/pagination"
	"github.com/haiyiyun/utils/http/request"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/utils/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (self *Service) Route_GET_MyContents(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var requestSFCL predefined.RequestServeMyFollowContentList
	if err := validator.FormStruct(&requestSFCL, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	fcModel := follow_content.NewModel(self.M)
	filter := fcModel.FilterByUserID(claims.UserID)

	if len(requestSFCL.Types) > 0 {
		filter = append(filter, fcModel.FilterByTypes(requestSFCL.Types)...)
	}

	if requestSFCL.OnlyUnreaded {
		filter = append(filter, bson.D{
			{"readed_time", bson.D{
				{"$eq", time.Time{}},
			}},
		}...)
	} else if requestSFCL.OnlyReaded {
		filter = append(filter, bson.D{
			{"readed_time", bson.D{
				{"$ne", time.Time{}},
			}},
		}...)
	} else {
		filterReadedTime := bson.D{
			{"readed_time", bson.D{
				{"$ne", time.Time{}},
			}},
		}
		if !requestSFCL.GteReadedTime.Time.IsZero() {
			filterReadedTime = append(filterReadedTime, bson.D{
				{"readed_time", bson.D{
					{"$gte", requestSFCL.GteReadedTime.Time},
				}},
			}...)
		}
		if !requestSFCL.LteReadedTime.Time.IsZero() {
			filterReadedTime = append(filterReadedTime, bson.D{
				{"readed_time", bson.D{
					{"$lte", requestSFCL.LteReadedTime.Time},
				}},
			}...)
		}

		if len(filterReadedTime) > 1 {
			filter = append(filter, bson.D{
				{"$or", []bson.D{
					{{"readed_time", bson.D{
						{"$eq", time.Time{}},
					}}},
					filterReadedTime,
				}},
			}...)
		} else {
			filter = append(filter, bson.D{
				{"readed_time", bson.D{
					{"$eq", time.Time{}},
				}},
			}...)
		}

	}

	cnt, _ := fcModel.CountDocuments(r.Context(), filter)
	pg := pagination.Parse(r, cnt)

	projection := bson.D{
		{"_id", 1},
		{"follow_relationship_id", 1},
		{"type", 1},
		{"user_id", 1},
		{"content_id", 1},
		{"readed_time", 1},
		{"extension_id", 1},
		{"create_time", 1},
	}

	opt := options.Find().SetSort(bson.D{
		{"create_time", -1},
	}).SetProjection(projection).SetSkip(pg.SkipNum).SetLimit(pg.PageSize)

	if cur, err := fcModel.Find(r.Context(), filter, opt); err != nil {
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
