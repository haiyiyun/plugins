package content

import (
	"net/http"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/mongodb/geometry"
	"github.com/haiyiyun/plugins/content/database/model"
	"github.com/haiyiyun/plugins/content/database/model/content"
	"github.com/haiyiyun/plugins/content/predefined"
	"github.com/haiyiyun/utils/help"
	"github.com/haiyiyun/utils/http/pagination"
	"github.com/haiyiyun/utils/http/request"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/utils/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	var requestCC predefined.RequestServeContentCreate
	if err := validator.FormStruct(&requestCC, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	coordinates := geometry.PointCoordinates{
		requestCC.Longitude, requestCC.Latitude,
	}

	link := model.ContentLink{
		URL: requestCC.ContentLink,
	}

	tags := []string{}
	tags = append(tags, requestCC.UserTags...)

	//TODO 发布干预
	status := predefined.PublishStatusNormal

	//TODO 伪装处理
	guise := model.ContentGuise{}

	contentModel := content.NewModel(self.M)

	//关联类型处理
	if requestCC.AssociateID != primitive.NilObjectID {
		filter := contentModel.FilterNormalContent()
		filter = append(filter, contentModel.FilterByID(requestCC.AssociateID)...)
		if sr := contentModel.FindOne(r.Context(), filter, options.FindOne().SetProjection(bson.D{
			{"type", 1},
			{"publish_type", 1},
			{"limit_associate_type", 1},
			{"limit_associate_num", 1},
		})); sr.Err() != nil {
			log.Error(sr.Err())
			response.JSON(rw, http.StatusBadRequest, nil, "400404")
			return
		} else {
			var cont model.Content
			if err := sr.Decode(&cont); err != nil {
				log.Error(err)
				response.JSON(rw, http.StatusBadRequest, nil, "400000")
				return
			} else {
				//判断限制关联类型
				if cont.LimitAssociateType > 0 {
					if cont.LimitAssociateType != requestCC.PublishType {
						response.JSON(rw, http.StatusForbidden, nil, "403010")
						return
					}
				}

				//判断限制关联数量
				if cont.LimitAssociateNum > 0 {
					filterLimit := contentModel.FilterNormalContent()
					filterLimit = append(filterLimit, bson.D{
						{"associate_id", requestCC.AssociateID},
					}...)
					if cnt, err := contentModel.CountDocuments(r.Context(), filterLimit); err != nil && err != mongo.ErrNoDocuments {
						log.Error(err)
						response.JSON(rw, http.StatusBadRequest, nil, "400010")
						return
					} else {
						if cnt >= int64(cont.LimitAssociateNum) {
							response.JSON(rw, http.StatusForbidden, nil, "403020")
							return
						}
					}
				}
			}
		}
	}

	ctnt := &model.Content{
		PublishUserID:                          userID,
		Type:                                   requestCC.Type,
		PublishType:                            requestCC.PublishType,
		AssociateType:                          requestCC.AssociateType,
		AssociateID:                            requestCC.AssociateID,
		LimitAssociateType:                     requestCC.LimitAssociateType,
		LimitAllDiscussNum:                     requestCC.LimitAssociateNum,
		CategoryID:                             requestCC.CategoryID,
		SubjectID:                              requestCC.SubjectID,
		AtUsers:                                requestCC.AtUsers,
		Author:                                 requestCC.Author,
		Title:                                  requestCC.Title,
		Cover:                                  requestCC.Cover,
		Description:                            requestCC.Description,
		Video:                                  requestCC.Video,
		Voice:                                  requestCC.Voice,
		Images:                                 requestCC.Images,
		Content:                                requestCC.Content,
		Link:                                   link,
		Location:                               geometry.NewPoint(coordinates),
		UserTags:                               requestCC.UserTags,
		Visibility:                             requestCC.Visibility,
		Copy:                                   requestCC.Copy,
		OnlyUserIDDiscuss:                      requestCC.OnlyUserIDDiscuss,
		OnlyUserIDCanReplyDiscuss:              requestCC.OnlyUserIDCanReplyDiscuss,
		OnlyUserIDCanNotReplyDiscuss:           requestCC.OnlyUserIDCanNotReplyDiscuss,
		LimitAssociateNum:                      requestCC.LimitAllDiscussNum,
		LimitPublishUserDiscussNum:             requestCC.LimitPublishUserDiscussNum,
		LimitNotPublishUserAllUserDiscussNum:   requestCC.LimitNotPublishUserAllUserDiscussNum,
		LimitNotPublishUserEveryUserDiscussNum: requestCC.LimitNotPublishUserEveryUserDiscussNum,
		ForbidForward:                          requestCC.ForbidForward,
		ForbidDownload:                         requestCC.ForbidDownload,
		ForbidDiscuss:                          requestCC.ForbidDiscuss,
		Tags:                                   tags,
		Guise:                                  guise,
		Status:                                 status,
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

	var requestCL predefined.RequestServeContentList
	if err := validator.FormStruct(&requestCL, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	coordinates := geometry.PointCoordinates{
		requestCL.Longitude, requestCL.Latitude,
	}

	contentModel := content.NewModel(self.M)
	filter := contentModel.FilterNormalContent()
	filter = append(filter, contentModel.FilterByTypes(requestCL.Types)...)

	if requestCL.Visibility == predefined.VisibilityTypeSelf {
		filter = append(filter, contentModel.FilterByPublishUserID(userID)...)
		filter = append(filter, contentModel.FilterByVisibility(requestCL.Visibility)...)
	} else {
		if requestCL.PublishUserID != primitive.NilObjectID {
			filter = append(filter, contentModel.FilterByPublishUserID(requestCL.PublishUserID)...)
			filter = append(filter, contentModel.FilterByVisibility(requestCL.Visibility)...)
		} else {
			if requestCL.Visibility != predefined.VisibilityTypeAll {
				filter = append(filter, contentModel.FilterByVisibilityOrAll(requestCL.Visibility)...)
			} else {
				filter = append(filter, contentModel.FilterByVisibility(requestCL.Visibility)...)
			}
		}
	}

	if requestCL.CategoryID != primitive.NilObjectID {
		filter = append(filter, contentModel.FilterByCategoryID(requestCL.CategoryID)...)
	}

	if requestCL.SubjectID != primitive.NilObjectID {
		filter = append(filter, contentModel.FilterBySubjectID(requestCL.SubjectID)...)
	}

	if requestCL.DiscussTotalGte > 0 {
		filter = append(filter, contentModel.FilterByGteDiscussEstimateTotal(requestCL.DiscussTotalGte)...)
	}

	if requestCL.DiscussTotalLte > 0 {
		filter = append(filter, contentModel.FilterByLteDiscussEstimateTotal(requestCL.DiscussTotalLte)...)
	}

	if requestCL.DiscussTotalZero {
		filter = append(filter, contentModel.FilterByDiscussEstimateTotal(0)...)
	}

	if len(requestCL.Tags) > 0 {
		filter = append(filter, bson.E{
			"tags", bson.D{
				{"$in", requestCL.Tags},
			},
		})
	}

	if coordinates != geometry.NilPointCoordinates {
		filter = append(filter, contentModel.FilterByLocation(geometry.NewPoint(coordinates), requestCL.MaxDistance, requestCL.MinDistance)...)
	}

	cnt, _ := contentModel.CountDocuments(r.Context(), filter)
	pg := pagination.Parse(r, cnt)

	projection := bson.D{
		{"_id", 1},
		{"publish_user_id", 1},
		{"type", 1},
		{"publish_type", 1},
		{"associate_type", 1},
		{"associate_id", 1},
		{"category_id", 1},
		{"subject_id", 1},
		{"author", 1},
		{"title", 1},
		{"cover", 1},
		{"description", 1},
		{"copy", 1},
		{"bestest", 1},
		{"reliable", 1},
		{"guise", 1},
		{"anti_guise_user", 1},
		{"status", 1},
		{"discuss_estimate_total", 1},
		{"create_time", 1},
		{"update_time", 1},
	}

	opt := options.Find().SetSort(bson.D{
		{"create_time", -1},
	}).SetProjection(projection).SetSkip(pg.SkipNum).SetLimit(pg.PageSize)

	if cur, err := contentModel.Find(r.Context(), filter, opt); err != nil {
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

func (self *Service) Route_GET_Detail(rw http.ResponseWriter, r *http.Request) {
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

	var requestCD predefined.RequestServeContentDetail
	if err := validator.FormStruct(&requestCD, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	contentModel := content.NewModel(self.M)
	filter := contentModel.FilterByID(requestCD.ID)
	filter = append(filter, contentModel.FilterNormalContent()...)
	if sr := contentModel.FindOne(r.Context(), filter); sr.Err() == nil {
		var contentDetail model.Content
		if err := sr.Decode(&contentDetail); err == nil {
			if !contentDetail.HideDetail ||
				(contentDetail.HideDetail &&
					len(contentDetail.OnlyUserIDShowDetail) > 0 &&
					help.NewSlice(help.NewSlice(contentDetail.OnlyUserIDShowDetail).ObjectIDToStrings()).CheckItem(userID.Hex())) {
				response.JSON(rw, 0, contentDetail, "")
			} else {
				response.JSON(rw, http.StatusForbidden, nil, "")
			}
		} else {
			log.Error(err)
			response.JSON(rw, http.StatusServiceUnavailable, nil, "")
		}
	} else {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	}
}
