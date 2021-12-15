package discuss

import (
	"net/http"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/mongodb/geometry"
	"github.com/haiyiyun/plugins/content/database/model"
	"github.com/haiyiyun/plugins/content/database/model/content"
	"github.com/haiyiyun/plugins/content/database/model/discuss"
	"github.com/haiyiyun/plugins/content/predefined"
	"github.com/haiyiyun/utils/help"
	"github.com/haiyiyun/utils/http/pagination"
	"github.com/haiyiyun/utils/http/request"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/utils/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	var requestDC predefined.RequestServeDiscussCreate
	if err := validator.FormStruct(&requestDC, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	coordinates := geometry.PointCoordinates{
		requestDC.Longitude, requestDC.Latitude,
	}

	//TODO 发布干预
	status := predefined.PublishStatusNormal

	//判断对应type的object_id是否存在
	contentModel := content.NewModel(self.M)
	switch requestDC.Type {
	case predefined.DiscussTypeDynamic:
		contentType := predefined.ContentPublishTypeDynamic
		filter := contentModel.FilterNormalContent()
		filter = append(filter, contentModel.FilterByID(requestDC.ObjectID)...)
		filter = append(filter, contentModel.FilterByPublishType(contentType)...)
		if cnt, err := contentModel.CountDocuments(r.Context(), filter); err != nil || cnt == 0 {
			log.Error(err)
			response.JSON(rw, http.StatusNotFound, nil, "40400")
			return
		}
	case predefined.DiscussTypeArticle:
		contentType := predefined.ContentPublishTypeArticle
		filter := contentModel.FilterNormalContent()
		filter = append(filter, contentModel.FilterByID(requestDC.ObjectID)...)
		filter = append(filter, contentModel.FilterByPublishType(contentType)...)
		if cnt, err := contentModel.CountDocuments(r.Context(), filter); err != nil || cnt == 0 {
			log.Error(err)
			response.JSON(rw, http.StatusNotFound, nil, "40401")
			return
		}
	case predefined.DiscussTypeQuestion:
		contentType := predefined.ContentPublishTypeQuestion
		filter := contentModel.FilterNormalContent()
		filter = append(filter, contentModel.FilterByID(requestDC.ObjectID)...)
		filter = append(filter, contentModel.FilterByPublishType(contentType)...)
		if cnt, err := contentModel.CountDocuments(r.Context(), filter); err != nil || cnt == 0 {
			log.Error(err)
			response.JSON(rw, http.StatusNotFound, nil, "40402")
			return
		}
	case predefined.DiscussTypeAnswer:
		contentType := predefined.ContentPublishTypeAnswer
		filter := contentModel.FilterNormalContent()
		filter = append(filter, contentModel.FilterByID(requestDC.ObjectID)...)
		filter = append(filter, contentModel.FilterByPublishType(contentType)...)
		if cnt, err := contentModel.CountDocuments(r.Context(), filter); err != nil || cnt == 0 {
			log.Error(err)
			response.JSON(rw, http.StatusNotFound, nil, "40403")
			return
		}
	}

	discussModel := discuss.NewModel(self.M)
	dis := &model.Discuss{
		Type:           requestDC.Type,
		ObjectID:       requestDC.ObjectID,
		PublishUserID:  userID,
		AtUsers:        requestDC.AtUsers,
		ReplyDiscussID: requestDC.ReplyDiscussID,
		Text:           requestDC.Text,
		Location:       geometry.NewPoint(coordinates),
		Visibility:     requestDC.Visibility,
		Status:         status,
	}

	if ior, err := discussModel.Create(r.Context(), dis); err != nil || ior.InsertedID == nil {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, ior.InsertedID, "")
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

	var requestDL predefined.RequestServeDiscussList
	if err := validator.FormStruct(&requestDL, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	coordinates := geometry.PointCoordinates{
		requestDL.Longitude, requestDL.Latitude,
	}

	discussModel := discuss.NewModel(self.M)
	filter := discussModel.FilterNormalDiscuss()
	filter = append(filter, discussModel.FilterByType(requestDL.Type)...)
	filter = append(filter, discussModel.FilterByObjectID(requestDL.ObjectID)...)

	if requestDL.Visibility == predefined.VisibilityTypeSelf {
		filter = append(filter, discussModel.FilterByPublishUserID(userID)...)
		filter = append(filter, discussModel.FilterByVisibility(requestDL.Visibility)...)
	} else {
		if requestDL.PublishUserID != primitive.NilObjectID {
			filter = append(filter, discussModel.FilterByPublishUserID(requestDL.PublishUserID)...)
			filter = append(filter, discussModel.FilterByVisibility(requestDL.Visibility)...)
		} else {
			filter = append(filter, discussModel.FilterByVisibilityOrAll(requestDL.Visibility)...)
		}
	}

	if coordinates != geometry.NilPointCoordinates {
		filter = append(filter, discussModel.FilterByLocation(geometry.NewPoint(coordinates), requestDL.MaxDistance, requestDL.MinDistance)...)
	}

	cnt, _ := discussModel.CountDocuments(r.Context(), filter)
	pg := pagination.Parse(r, cnt)

	projection := bson.D{}

	opt := options.Find().SetSort(bson.D{
		{"create_time", -1},
	}).SetProjection(projection).SetSkip(pg.SkipNum).SetLimit(pg.PageSize)

	if cur, err := discussModel.Find(r.Context(), filter, opt); err != nil {
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
