package content

import (
	"net/http"
	"strconv"

	"github.com/haiyiyun/mongodb/geometry"
	"github.com/haiyiyun/plugins/content/database/model"
	"github.com/haiyiyun/plugins/content/database/model/content"
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
	publishTypeStr := r.FormValue("publish_type")
	publishType, _ := strconv.Atoi(publishTypeStr)
	associateTypeStr := r.FormValue("associate_type")
	associateType, _ := strconv.Atoi(associateTypeStr)
	associateIDStr := r.FormValue("associate_id")
	associateID, _ := primitive.ObjectIDFromHex(associateIDStr)
	categoryIDStr := r.FormValue("category_id")
	categoryID, _ := primitive.ObjectIDFromHex(categoryIDStr)
	subjectIDStr := r.FormValue("subject_id")
	subjectID, _ := primitive.ObjectIDFromHex(subjectIDStr)
	atUsersStrs := r.Form["at_users[]"]
	atUsers := help.NewSlice(atUsersStrs).ConvObjectID()
	author := r.FormValue("author")
	title := r.FormValue("title")
	cover := r.FormValue("cover")
	description := r.FormValue("description")

	video := r.FormValue("video")
	voice := r.FormValue("voice")
	images := r.Form["images[]"]
	contentStr := r.FormValue("content")
	contentLink := r.FormValue("content_link")
	link := model.ContentLink{
		URL: contentLink,
	}

	longitudeStr := r.FormValue("longitude") //经度
	latitudeStr := r.FormValue("latitude")   //维度
	longitude, _ := strconv.ParseFloat(longitudeStr, 64)
	latitude, _ := strconv.ParseFloat(latitudeStr, 64)
	coordinates := geometry.PointCoordinates{
		longitude, latitude,
	}

	userTags := r.Form["user_tags[]"]

	visibilityStr := r.FormValue("visibility")
	visibility, _ := strconv.Atoi(visibilityStr)

	copyContentStr := r.FormValue("copy")
	copyContent, _ := strconv.ParseBool(copyContentStr)
	forbidForwardStr := r.FormValue("forbid_forward")
	forbidForward, _ := strconv.ParseBool(forbidForwardStr)
	forbidDownloadStr := r.FormValue("forbid_download")
	forbidDownload, _ := strconv.ParseBool(forbidDownloadStr)
	forbidDiscussStr := r.FormValue("forbid_discuss")
	forbidDiscuss, _ := strconv.ParseBool(forbidDiscussStr)

	valid := validator.Validation{}
	valid.Digital(typeStr).Key("type_str").Message("type必须数字")
	valid.Have(
		typ,
		predefined.ContentTypeVideoShort,
		predefined.ContentTypeVideoLong,
		predefined.ContentTypeVoiceShort,
		predefined.ContentTypeVoiceLong,
		predefined.ContentTypeImageText,
		predefined.ContentTypeImage,
		predefined.ContentTypeText,
	).Key("type").Message("type必须支持的类型")

	valid.Digital(publishTypeStr).Key("publish_type_str").Message("publish_type必须数字")

	valid.Have(
		publishType,
		predefined.ContentPublishTypeDynamic,
		predefined.ContentPublishTypeArticle,
		predefined.ContentPublishTypeQuestion,
		predefined.ContentPublishTypeAnswer,
	).Key("type").Message("type必须支持的类型")

	valid.Digital(associateTypeStr).Key("associate_type_str").Message("associate_type必须数字")

	valid.Have(
		associateType,
		predefined.ContentAssociateTypeSelf,
		predefined.ContentAssociateTypeForward,
		predefined.ContentAssociateTypeCollection,
		predefined.ContentAssociateTypeDynamic,
		predefined.ContentAssociateTypeArticle,
		predefined.ContentAssociateTypeQuestion,
		predefined.ContentAssociateTypeAnswer,
	).Key("associate_type").Message("associate_type必须支持的类型")

	if associateType != predefined.ContentAssociateTypeSelf {
		valid.Required(associateIDStr).Key("associate_id").Message("associate_id不能为空")
	}

	valid.Required(title).Key("title").Message("title不能为空")
	valid.Required(cover).Key("cover").Message("cover不能为空")

	switch typ {
	case predefined.ContentTypeVideoShort:
		valid.Required(video).Key("video").Message("video不能为空")
	case predefined.ContentTypeVideoLong:
		valid.Required(video).Key("video").Message("video不能为空")
	case predefined.ContentTypeVoiceShort:
		valid.Required(voice).Key("voice").Message("voice不能为空")
	case predefined.ContentTypeVoiceLong:
		valid.Required(voice).Key("voice").Message("voice不能为空")
	case predefined.ContentTypeImageText:
		valid.Required(images).Key("images[]").Message("images[]不能为空")
		valid.Required(contentStr).Key("content").Message("content不能为空")
	case predefined.ContentTypeImage:
		valid.Required(associateIDStr).Key("images[]").Message("images[]不能为空")
	case predefined.ContentTypeText:
		valid.Required(contentStr).Key("content").Message("content不能为空")
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

	tags := []string{}
	tags = append(tags, userTags...)

	contentModel := content.NewModel(self.M)
	ctnt := &model.Content{
		PublishUserID:  userID,
		Type:           typ,
		PublishType:    publishType,
		AssociateType:  associateType,
		AssociateID:    associateID,
		CategoryID:     categoryID,
		SubjectID:      subjectID,
		AtUsers:        atUsers,
		Author:         author,
		Title:          title,
		Cover:          cover,
		Description:    description,
		Video:          video,
		Voice:          voice,
		Images:         images,
		Content:        contentStr,
		Link:           link,
		Location:       geometry.NewPoint(coordinates),
		UserTags:       userTags,
		Visibility:     visibility,
		Copy:           copyContent,
		ForbidForward:  forbidForward,
		ForbidDownload: forbidDownload,
		ForbidDiscuss:  forbidDiscuss,
		Tags:           tags,
	}

	if ior, err := contentModel.Create(r.Context(), ctnt); err != nil || ior.InsertedID == nil {
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
	visibilityStr := r.URL.Query().Get("visibility")
	visibility, _ := strconv.Atoi(visibilityStr)
	tags := r.URL.Query()["tags[]"]
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
		predefined.ContentTypeVideoShort,
		predefined.ContentTypeVideoLong,
		predefined.ContentTypeVoiceShort,
		predefined.ContentTypeVoiceLong,
		predefined.ContentTypeImageText,
		predefined.ContentTypeImage,
		predefined.ContentTypeText,
		predefined.ContentTypeSystemVideoShort,
		predefined.ContentTypeSystemVideoLong,
		predefined.ContentTypeSystemVoiceShort,
		predefined.ContentTypeSystemVoiceLong,
		predefined.ContentTypeSystemImageText,
		predefined.ContentTypeSystemImage,
		predefined.ContentTypeSystemText,
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

	if publishUserIDStr != "" {
		valid.BsonObjectID(publishUserIDStr).Key("publish_user_id").Message("publish_user_id必须支持的格式")
	}

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	contentModel := content.NewModel(self.M)
	filter := contentModel.FilterNormalContent()
	filter = append(filter, contentModel.FilterByType(typ)...)

	if visibility == predefined.VisibilityTypeSelf {
		filter = append(filter, contentModel.FilterByPublishUserID(userID)...)
		filter = append(filter, contentModel.FilterByVisibility(visibility)...)
	} else {
		if publishUserID != primitive.NilObjectID {
			filter = append(filter, contentModel.FilterByPublishUserID(publishUserID)...)
			filter = append(filter, contentModel.FilterByVisibility(visibility)...)
		} else {
			filter = append(filter, contentModel.FilterByVisibilityOrAll(visibility)...)
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
		filter = append(filter, contentModel.FilterByLocation(geometry.NewPoint(coordinates), maxDistance, minDistance)...)
	}

	cnt, _ := contentModel.CountDocuments(r.Context(), filter)
	pg := pagination.Parse(r, cnt)

	opt := options.Find().SetSort(bson.D{
		{"update_time", -1},
	}).SetProjection(bson.D{}).SetSkip(pg.SkipNum).SetLimit(pg.PageSize)

	if cur, err := contentModel.Find(r.Context(), filter, opt); err != nil {
		response.JSON(rw, http.StatusNotFound, nil, "")
	} else {
		items := []model.Content{}
		if err := cur.All(r.Context(), &items); err != nil {
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
