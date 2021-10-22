package subject

import (
	"net/http"
	"strconv"

	"github.com/haiyiyun/plugins/content/database/model"
	"github.com/haiyiyun/plugins/content/database/model/subject"
	"github.com/haiyiyun/plugins/content/predefined"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/haiyiyun/mongodb/geometry"
	"github.com/haiyiyun/utils/http/request"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/validator"
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
	subjectStr := r.FormValue("subject")
	visibilityStr := r.FormValue("visibility")
	visibility, _ := strconv.Atoi(visibilityStr)

	longitudeStr := r.FormValue("longitude") //经度
	latitudeStr := r.FormValue("latitude")   //维度
	longitude, _ := strconv.ParseFloat(longitudeStr, 64)
	latitude, _ := strconv.ParseFloat(latitudeStr, 64)
	coordinates := geometry.PointCoordinates{
		longitude, latitude,
	}

	valid := validator.Validation{}
	valid.Digital(typeStr).Key("type_str").Message("type必须数字")
	valid.Have(typ,
		predefined.SubjectTypeUserDynamic,
		predefined.SubjectTypeUserArticle,
	).Key("type").Message("type必须是支持的类型")

	valid.Required(subjectStr).Key("subject").Message("subject不能为空")

	valid.Digital(visibilityStr).Key("visibility_str").Message("visibility必须数字")
	valid.Have(visibility,
		predefined.VisibilityTypeSelf,
		predefined.VisibilityTypeHome,
		predefined.VisibilityTypeRelationship,
		predefined.VisibilityTypeStranger,
		predefined.VisibilityTypeSubject,
		predefined.VisibilityTypeNearly,
	).Key("visibility").Message("visibility必须是支持的类型")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	subjectModel := subject.NewModel(self.M)
	if _, err := subjectModel.Create(r.Context(), model.Subject{
		PublishUserID: userID,
		Type:          typ,
		Subject:       subjectStr,
		Enable:        true,
		Location:      geometry.NewPoint(coordinates),
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

	typeStr := r.URL.Query().Get("type")
	typ, _ := strconv.Atoi(typeStr)
	visibilityStr := r.URL.Query().Get("visibility")
	visibility, _ := strconv.Atoi(visibilityStr)
	tags := r.URL.Query()["tags[]"]
	publishUserIdStr := r.URL.Query().Get("publish_user_id")
	publishUserId, _ := primitive.ObjectIDFromHex(publishUserIdStr)

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
		predefined.SubjectTypeSystemDynamic,
		predefined.SubjectTypeSystemArticle,
		predefined.SubjectTypeUserDynamic,
		predefined.SubjectTypeUserArticle,
	).Key("type").Message("type必须是支持的类型")

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

	if publishUserIdStr != "" {
		valid.BsonObjectID(publishUserIdStr).Key("publish_user_id").Message("publish_user_id必须支持的格式")
	}

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	subjectModel := subject.NewModel(self.M)
	filter := subjectModel.FilterNormalSubject()
	filter = append(filter, subjectModel.FilterByType(typ)...)

	if visibility == predefined.VisibilityTypeSelf {
		filter = append(filter, subjectModel.FilterByPublishUserID(userID)...)
		filter = append(filter, subjectModel.FilterByVisibility(visibility)...)
	} else {
		if publishUserId != primitive.NilObjectID {
			filter = append(filter, subjectModel.FilterByPublishUserID(publishUserId)...)
			filter = append(filter, subjectModel.FilterByVisibility(visibility)...)
		} else {
			filter = append(filter, subjectModel.FilterByVisibilityOrAll(visibility)...)
		}
	}

	if len(tags) > 0 {
		filter = append(filter, bson.E{
			"tags", bson.D{
				{"$in", tags},
			},
		})
	}

	if coordinates != geometry.NilPointCoordinates {
		filter = append(filter, subjectModel.FilterByLocation(geometry.NewPoint(coordinates), maxDistance, minDistance)...)
	}

	if cur, err := subjectModel.Find(r.Context(), filter); err != nil {
		response.JSON(rw, http.StatusNotFound, nil, "")
	} else {
		subjects := []model.Subject{}
		if err := cur.All(r.Context(), &subjects); err != nil {
			response.JSON(rw, http.StatusServiceUnavailable, nil, "")
		} else {
			response.JSON(rw, 0, subjects, "")
		}
	}
}
