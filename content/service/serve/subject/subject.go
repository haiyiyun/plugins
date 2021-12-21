package subject

import (
	"net/http"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/plugins/content/database/model"
	"github.com/haiyiyun/plugins/content/database/model/subject"
	"github.com/haiyiyun/plugins/content/predefined"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/haiyiyun/mongodb/geometry"
	"github.com/haiyiyun/utils/help"
	"github.com/haiyiyun/utils/http/pagination"
	"github.com/haiyiyun/utils/http/request"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/utils/validator"
)

func (self *Service) Route_POST_Create(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	userID, _ := primitive.ObjectIDFromHex(claims.Audience)

	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestSC predefined.RequestServeSubjectCreate
	if err := validator.FormStruct(&requestSC, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	coordinates := geometry.PointCoordinates{
		requestSC.Longitude, requestSC.Latitude,
	}

	tags := []string{}
	tags = append(tags, requestSC.UserTags...)

	//TODO 发布干预
	status := predefined.PublishStatusNormal

	subjectModel := subject.NewModel(self.M)
	if _, err := subjectModel.Create(r.Context(), model.Subject{
		PublishUserID: userID,
		Type:          requestSC.Type,
		Subject:       requestSC.Subject,
		Enable:        true,
		UserTags:      requestSC.UserTags,
		Tags:          tags,
		Location:      geometry.NewPoint(coordinates),
		Status:        status,
	}); err != nil {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_GET_List(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	userID, _ := primitive.ObjectIDFromHex(claims.Audience)

	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var requestSL predefined.RequestServeSubjectList
	if err := validator.FormStruct(&requestSL, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	coordinates := geometry.PointCoordinates{
		requestSL.Longitude, requestSL.Latitude,
	}

	subjectModel := subject.NewModel(self.M)
	filter := subjectModel.FilterNormalSubject()
	filter = append(filter, subjectModel.FilterByType(requestSL.Type)...)

	if requestSL.Visibility == predefined.VisibilityTypeSelf {
		filter = append(filter, subjectModel.FilterByPublishUserID(userID)...)
		filter = append(filter, subjectModel.FilterByVisibility(requestSL.Visibility)...)
	} else {
		if requestSL.PublishUserID != primitive.NilObjectID {
			filter = append(filter, subjectModel.FilterByPublishUserID(requestSL.PublishUserID)...)
			filter = append(filter, subjectModel.FilterByVisibility(requestSL.Visibility)...)
		} else {
			if requestSL.Visibility != predefined.VisibilityTypeAll {
				filter = append(filter, subjectModel.FilterByVisibilityOrAll(requestSL.Visibility)...)
			} else {
				filter = append(filter, subjectModel.FilterByVisibility(requestSL.Visibility)...)
			}
		}
	}

	if len(requestSL.Tags) > 0 {
		filter = append(filter, bson.E{
			"tags", bson.D{
				{"$in", requestSL.Tags},
			},
		})
	}

	if coordinates != geometry.NilPointCoordinates {
		filter = append(filter, subjectModel.FilterByLocation(geometry.NewPoint(coordinates), requestSL.MaxDistance, requestSL.MinDistance)...)
	}

	cnt, _ := subjectModel.CountDocuments(r.Context(), filter)
	pg := pagination.Parse(r, cnt)

	projection := bson.D{
		{"_id", 1},
		{"publish_user_id", 1},
		{"type", 1},
		{"subject", 1},
		{"cover", 1},
		{"description", 1},
		{"status", 1},
	}

	opt := options.Find().SetSort(bson.D{
		{"create_time", -1},
	}).SetProjection(projection).SetSkip(pg.SkipNum).SetLimit(pg.PageSize)

	if cur, err := subjectModel.Find(r.Context(), filter, opt); err != nil {
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
