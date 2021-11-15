package discuss

import (
	"net/http"
	"strconv"

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
	"github.com/haiyiyun/validator"
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
	typeStr := r.FormValue("type")
	typ, _ := strconv.Atoi(typeStr)
	objectIDStr := r.FormValue("object_id")
	objectID, _ := primitive.ObjectIDFromHex(objectIDStr)
	atUsersStrs := r.Form["at_users[]"]
	atUsers := help.NewSlice(atUsersStrs).ConvObjectID()
	replyDiscussIDStr := r.FormValue("reply_discuss_id")
	replyDiscussID, _ := primitive.ObjectIDFromHex(replyDiscussIDStr)
	text := r.FormValue("text")

	longitudeStr := r.FormValue("longitude") //经度
	latitudeStr := r.FormValue("latitude")   //维度
	longitude, _ := strconv.ParseFloat(longitudeStr, 64)
	latitude, _ := strconv.ParseFloat(latitudeStr, 64)
	coordinates := geometry.PointCoordinates{
		longitude, latitude,
	}

	visibilityStr := r.FormValue("visibility")
	visibility, _ := strconv.Atoi(visibilityStr)

	valid := validator.Validation{}
	valid.Digital(typeStr).Key("type_str").Message("type必须数字")
	valid.Have(
		typ,
		predefined.DiscussTypeDynamic,
		predefined.DiscussTypeArticle,
		predefined.DiscussTypeQuestion,
		predefined.DiscussTypeAnswer,
	).Key("type").Message("type必须支持的类型")

	valid.BsonObjectID(objectIDStr).Key("object_id").Message("object_id必须支持的格式")

	if replyDiscussIDStr != "" {
		valid.BsonObjectID(replyDiscussIDStr).Key("reply_discuss_id").Message("reply_discuss_id必须支持的格式")
	}

	valid.Digital(visibilityStr).Key("visibility_str").Message("visibility必须数字")
	valid.Have(visibility,
		predefined.VisibilityTypeSelf,
		predefined.VisibilityTypeHome,
		predefined.VisibilityTypeRelationship,
		predefined.VisibilityTypeStranger,
		predefined.VisibilityTypeSubject,
		predefined.VisibilityTypeNearly,
		predefined.VisibilityTypeCity,
		predefined.VisibilityTypeProvince,
		predefined.VisibilityTypeNation,
		predefined.VisibilityTypeAll,
	).Key("visibility").Message("visibility必须是支持的类型")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().Message)
		return
	}

	//TODO 发布干预
	status := predefined.PublishStatusNormal

	//判断对应type的object_id是否存在
	contentModel := content.NewModel(self.M)
	switch typ {
	case predefined.DiscussTypeDynamic:
		contentType := predefined.ContentPublishTypeDynamic
		filter := contentModel.FilterNormalContent()
		filter = append(filter, contentModel.FilterByID(objectID)...)
		filter = append(filter, contentModel.FilterByPublishType(contentType)...)
		if cnt, err := contentModel.CountDocuments(r.Context(), filter); err != nil || cnt == 0 {
			log.Error(err)
			response.JSON(rw, http.StatusNotFound, nil, "40400")
			return
		}
	case predefined.DiscussTypeArticle:
		contentType := predefined.ContentPublishTypeArticle
		filter := contentModel.FilterNormalContent()
		filter = append(filter, contentModel.FilterByID(objectID)...)
		filter = append(filter, contentModel.FilterByPublishType(contentType)...)
		if cnt, err := contentModel.CountDocuments(r.Context(), filter); err != nil || cnt == 0 {
			log.Error(err)
			response.JSON(rw, http.StatusNotFound, nil, "40401")
			return
		}
	case predefined.DiscussTypeQuestion:
		contentType := predefined.ContentPublishTypeQuestion
		filter := contentModel.FilterNormalContent()
		filter = append(filter, contentModel.FilterByID(objectID)...)
		filter = append(filter, contentModel.FilterByPublishType(contentType)...)
		if cnt, err := contentModel.CountDocuments(r.Context(), filter); err != nil || cnt == 0 {
			log.Error(err)
			response.JSON(rw, http.StatusNotFound, nil, "40402")
			return
		}
	case predefined.DiscussTypeAnswer:
		contentType := predefined.ContentPublishTypeAnswer
		filter := contentModel.FilterNormalContent()
		filter = append(filter, contentModel.FilterByID(objectID)...)
		filter = append(filter, contentModel.FilterByPublishType(contentType)...)
		if cnt, err := contentModel.CountDocuments(r.Context(), filter); err != nil || cnt == 0 {
			log.Error(err)
			response.JSON(rw, http.StatusNotFound, nil, "40403")
			return
		}
	}

	discussModel := discuss.NewModel(self.M)
	dis := &model.Discuss{
		Type:           typ,
		ObjectID:       objectID,
		PublishUserID:  userID,
		AtUsers:        atUsers,
		ReplyDiscussID: replyDiscussID,
		Text:           text,
		Location:       geometry.NewPoint(coordinates),
		Visibility:     visibility,
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

	typeStr := r.URL.Query().Get("type")
	typ, _ := strconv.Atoi(typeStr)
	objectIDStr := r.URL.Query().Get("object_id")
	objectID, _ := primitive.ObjectIDFromHex(objectIDStr)
	visibilityStr := r.URL.Query().Get("visibility")
	visibility, _ := strconv.Atoi(visibilityStr)
	publishUserIDStr := r.URL.Query().Get("publish_user_id")
	publishUserID, _ := primitive.ObjectIDFromHex(publishUserIDStr)

	longitudeStr := r.URL.Query().Get("longitude") //经度
	latitudeStr := r.URL.Query().Get("latitude")   //维度
	longitude, _ := strconv.ParseFloat(longitudeStr, 64)
	latitude, _ := strconv.ParseFloat(latitudeStr, 64)
	coordinates := geometry.PointCoordinates{
		longitude, latitude,
	}

	maxDistanceStr := r.URL.Query().Get("max_distance")
	maxDistance, _ := strconv.ParseFloat(maxDistanceStr, 64)
	minDistanceStr := r.URL.Query().Get("min_distance")
	minDistance, _ := strconv.ParseFloat(minDistanceStr, 64)

	valid := validator.Validation{}
	valid.Digital(typeStr).Key("type_str").Message("type必须数字")
	valid.Have(typ,
		predefined.DiscussTypeDynamic,
		predefined.DiscussTypeArticle,
		predefined.DiscussTypeQuestion,
		predefined.DiscussTypeAnswer,
	).Key("type").Message("type必须是支持的类型")

	valid.BsonObjectID(objectIDStr).Key("object_id").Message("object_id必须支持的格式")

	valid.Digital(visibilityStr).Key("visibility_str").Message("visibility必须数字")
	valid.Have(visibility,
		predefined.VisibilityTypeSelf,
		predefined.VisibilityTypeHome,
		predefined.VisibilityTypeRelationship,
		predefined.VisibilityTypeStranger,
		predefined.VisibilityTypeSubject,
		predefined.VisibilityTypeNearly,
		predefined.VisibilityTypeCity,
		predefined.VisibilityTypeProvince,
		predefined.VisibilityTypeNation,
		predefined.VisibilityTypeAll,
	).Key("visibility").Message("visibility必须是支持的类型")

	if publishUserIDStr != "" {
		valid.BsonObjectID(publishUserIDStr).Key("publish_user_id").Message("publish_user_id必须支持的格式")
	}

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	discussModel := discuss.NewModel(self.M)
	filter := discussModel.FilterNormalDiscuss()
	filter = append(filter, discussModel.FilterByType(typ)...)
	filter = append(filter, discussModel.FilterByObjectID(objectID)...)

	if visibility == predefined.VisibilityTypeSelf {
		filter = append(filter, discussModel.FilterByPublishUserID(userID)...)
		filter = append(filter, discussModel.FilterByVisibility(visibility)...)
	} else {
		if publishUserID != primitive.NilObjectID {
			filter = append(filter, discussModel.FilterByPublishUserID(publishUserID)...)
			filter = append(filter, discussModel.FilterByVisibility(visibility)...)
		} else {
			filter = append(filter, discussModel.FilterByVisibilityOrAll(visibility)...)
		}
	}

	if coordinates != geometry.NilPointCoordinates {
		filter = append(filter, discussModel.FilterByLocation(geometry.NewPoint(coordinates), maxDistance, minDistance)...)
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
