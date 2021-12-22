package content

import (
	"net/http"
	"time"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/mongodb/geometry"
	"github.com/haiyiyun/plugins/content/database/model"
	"github.com/haiyiyun/plugins/content/database/model/category"
	"github.com/haiyiyun/plugins/content/database/model/content"
	"github.com/haiyiyun/plugins/content/database/model/subject"
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

	//判断Category是否存在
	if requestCC.CategoryID != primitive.NilObjectID {
		categoryModel := category.NewModel(self.M)
		categoryFilter := categoryModel.FilterByID(requestCC.CategoryID)
		categoryFilter = append(categoryFilter, categoryModel.FilterNormalCategory()...)
		if sr := categoryModel.FindOne(r.Context(), categoryFilter); sr.Err() != nil {
			log.Error(sr.Err())
			response.JSON(rw, http.StatusBadRequest, nil, "400414")
			return
		} else {
			var cate model.Category
			if err := sr.Decode(&cate); err != nil {
				log.Error(err)
				response.JSON(rw, http.StatusBadRequest, nil, "400020")
				return
			} else {
				//处理限制
				if cate.LimitUserAtLeastLevel > 0 {
					if claims.Level < cate.LimitUserAtLeastLevel {
						if !help.NewSlice(help.NewSlice(cate.OnlyUserIDNotLimitUserLevel).ObjectIDToStrings()).CheckItem(userID.Hex()) {
							response.JSON(rw, http.StatusForbidden, nil, "403030")
						}
					}
				}

				if len(cate.LimitUserRole) > 0 {
					foundRole := false
					for _, role := range cate.LimitUserRole {
						if foundRole = request.CheckUserRole(r, role); foundRole {
							break
						}
					}

					if !foundRole {
						if !help.NewSlice(help.NewSlice(cate.OnlyUserIDNotLimitUserRole).ObjectIDToStrings()).CheckItem(userID.Hex()) {
							response.JSON(rw, http.StatusForbidden, nil, "403040")
						}
					}
				}

				if len(cate.LimitUserTag) > 0 {
					foundTag := false
					for _, tag := range cate.LimitUserTag {
						if foundTag = request.CheckUserTag(r, tag); foundTag {
							break
						}
					}

					if !foundTag {
						if !help.NewSlice(help.NewSlice(cate.OnlyUserIDNotLimitUserTag).ObjectIDToStrings()).CheckItem(userID.Hex()) {
							response.JSON(rw, http.StatusForbidden, nil, "403050")
						}
					}
				}
			}
		}
	}

	//判断Subject是否存在
	if requestCC.SubjectID != primitive.NilObjectID {
		subjectModel := subject.NewModel(self.M)
		subjectFilter := subjectModel.FilterByID(requestCC.SubjectID)
		subjectFilter = append(subjectFilter, subjectModel.FilterNormalSubject()...)
		if sr := subjectModel.FindOne(r.Context(), subjectFilter); sr.Err() != nil {
			log.Error(sr.Err())
			response.JSON(rw, http.StatusBadRequest, nil, "400424")
			return
		} else {
			var subj model.Subject
			if err := sr.Decode(&subj); err != nil {
				log.Error(err)
				response.JSON(rw, http.StatusBadRequest, nil, "400030")
				return
			} else {
				//处理限制
				if subj.LimitUserAtLeastLevel > 0 {
					if claims.Level < subj.LimitUserAtLeastLevel {
						if !help.NewSlice(help.NewSlice(subj.OnlyUserIDNotLimitUserLevel).ObjectIDToStrings()).CheckItem(userID.Hex()) {
							response.JSON(rw, http.StatusForbidden, nil, "403060")
						}
					}
				}

				if len(subj.LimitUserRole) > 0 {
					foundRole := false
					for _, role := range subj.LimitUserRole {
						if foundRole = request.CheckUserRole(r, role); foundRole {
							break
						}
					}

					if !foundRole {
						if !help.NewSlice(help.NewSlice(subj.OnlyUserIDNotLimitUserRole).ObjectIDToStrings()).CheckItem(userID.Hex()) {
							response.JSON(rw, http.StatusForbidden, nil, "403070")
						}
					}
				}

				if len(subj.LimitUserTag) > 0 {
					foundTag := false
					for _, tag := range subj.LimitUserTag {
						if foundTag = request.CheckUserTag(r, tag); foundTag {
							break
						}
					}

					if !foundTag {
						if !help.NewSlice(help.NewSlice(subj.OnlyUserIDNotLimitUserTag).ObjectIDToStrings()).CheckItem(userID.Hex()) {
							response.JSON(rw, http.StatusForbidden, nil, "403080")
						}
					}
				}
			}
		}
	}

	if requestCC.AtUsers == nil {
		requestCC.AtUsers = []primitive.ObjectID{}
	}

	if requestCC.OnlyUserIDShowDetail == nil {
		requestCC.OnlyUserIDShowDetail = []primitive.ObjectID{}
	}

	if requestCC.OnlyUserIDDiscuss == nil {
		requestCC.OnlyUserIDDiscuss = []primitive.ObjectID{}
	}

	if requestCC.OnlyUserIDCanReplyDiscuss == nil {
		requestCC.OnlyUserIDCanReplyDiscuss = []primitive.ObjectID{}
	}

	if requestCC.OnlyUserIDCanNotReplyDiscuss == nil {
		requestCC.OnlyUserIDCanNotReplyDiscuss = []primitive.ObjectID{}
	}

	if requestCC.OnlyUserIDShowDiscuss == nil {
		requestCC.OnlyUserIDShowDiscuss = []primitive.ObjectID{}
	}

	if requestCC.OnlyUserIDNotLimitUserLevel == nil {
		requestCC.OnlyUserIDNotLimitUserLevel = []primitive.ObjectID{}
	}

	if requestCC.OnlyUserIDNotLimitUserRole == nil {
		requestCC.OnlyUserIDNotLimitUserRole = []primitive.ObjectID{}
	}

	if requestCC.OnlyUserIDNotLimitUserTag == nil {
		requestCC.OnlyUserIDNotLimitUserTag = []primitive.ObjectID{}
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
		Value:                                  requestCC.Value,
		HideDetail:                             requestCC.HideDetail,
		OnlyUserIDShowDetail:                   requestCC.OnlyUserIDShowDetail,
		Copy:                                   requestCC.Copy,
		LimitUserAtLeastLevel:                  requestCC.LimitUserAtLeastLevel,
		OnlyUserIDNotLimitUserLevel:            requestCC.OnlyUserIDNotLimitUserLevel,
		LimitUserRole:                          requestCC.LimitUserRole,
		OnlyUserIDNotLimitUserRole:             requestCC.OnlyUserIDNotLimitUserRole,
		LimitUserTag:                           requestCC.LimitUserTag,
		OnlyUserIDNotLimitUserTag:              requestCC.OnlyUserIDNotLimitUserTag,
		OnlyUserIDDiscuss:                      requestCC.OnlyUserIDDiscuss,
		OnlyUserIDCanReplyDiscuss:              requestCC.OnlyUserIDCanReplyDiscuss,
		OnlyUserIDCanNotReplyDiscuss:           requestCC.OnlyUserIDCanNotReplyDiscuss,
		LimitAssociateNum:                      requestCC.LimitAllDiscussNum,
		LimitPublishUserDiscussNum:             requestCC.LimitPublishUserDiscussNum,
		LimitNotPublishUserAllUserDiscussNum:   requestCC.LimitNotPublishUserAllUserDiscussNum,
		LimitNotPublishUserEveryUserDiscussNum: requestCC.LimitNotPublishUserEveryUserDiscussNum,
		HideDiscuss:                            requestCC.HideDiscuss,
		OnlyUserIDShowDiscuss:                  requestCC.OnlyUserIDShowDiscuss,
		ForbidForward:                          requestCC.ForbidForward,
		ForbidDownload:                         requestCC.ForbidDownload,
		ForbidDiscuss:                          requestCC.ForbidDiscuss,
		Tags:                                   tags,
		ReadedUser:                             []primitive.ObjectID{},
		WantedUser:                             []primitive.ObjectID{},
		LikedUser:                              []primitive.ObjectID{},
		HatedUser:                              []primitive.ObjectID{},
		Guise:                                  guise,
		AntiGuiseUser:                          []primitive.ObjectID{},
		StartTime:                              requestCC.StartTime.Time,
		EndTime:                                requestCC.EndTime.Time,
		ExtraData:                              requestCC.ExtraData,
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
		if len(requestCL.PublishUserID) > 0 {
			filter = append(filter, contentModel.FilterByPublishUserIDs(requestCL.PublishUserID)...)
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

	if requestCL.LimitUserAtLeastLevel > 0 {
		filter = append(filter, contentModel.FilterByGteLimitUserAtLeastLevel(requestCL.LimitUserAtLeastLevel)...)
	}

	if len(requestCL.LimitUserRole) > 0 {
		filter = append(filter, contentModel.FilterByLimitUserRole(requestCL.LimitUserRole)...)
	}

	if len(requestCL.LimitUserTag) > 0 {
		filter = append(filter, contentModel.FilterByLimitUserTag(requestCL.LimitUserTag)...)
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

	if requestCL.ValueGte > 0 {
		filter = append(filter, contentModel.FilterByGteValue(requestCL.ValueGte)...)
	}

	if requestCL.ValueLte > 0 {
		filter = append(filter, contentModel.FilterByLteValue(requestCL.ValueLte)...)
	}

	if requestCL.ValueZero {
		filter = append(filter, contentModel.FilterByValue(0)...)
	}

	if !requestCL.StartTime.Time.IsZero() {
		filter = append(filter, contentModel.FilterByGteStartTime(requestCL.StartTime.Time)...)
	}

	if !requestCL.EndTime.Time.IsZero() {
		filter = append(filter, contentModel.FilterByGteStartTime(requestCL.EndTime.Time)...)
	}

	if requestCL.InTime {
		now := time.Now()
		filter = append(filter, contentModel.FilterByGteStartTime(now)...)
		filter = append(filter, contentModel.FilterByGteStartTime(now)...)
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
		{"user_tags", 1},
		{"visibility", 1},
		{"value", 1},
		{"copy", 1},
		{"bestest", 1},
		{"reliable", 1},
		{"guise", 1},
		{"anti_guise_user", 1},
		{"start_time", 1},
		{"end_time", 1},
		{"extra_data", 1},
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
			now := time.Now()
			if contentDetail.LimitUserAtLeastLevel > 0 {
				if claims.Level < contentDetail.LimitUserAtLeastLevel {
					if !help.NewSlice(help.NewSlice(contentDetail.OnlyUserIDNotLimitUserLevel).ObjectIDToStrings()).CheckItem(userID.Hex()) {
						response.JSON(rw, http.StatusForbidden, nil, "403010")
					}
				}
			}

			if len(contentDetail.LimitUserRole) > 0 {
				foundRole := false
				for _, role := range contentDetail.LimitUserRole {
					if foundRole = request.CheckUserRole(r, role); foundRole {
						break
					}
				}

				if !foundRole {
					if !help.NewSlice(help.NewSlice(contentDetail.OnlyUserIDNotLimitUserRole).ObjectIDToStrings()).CheckItem(userID.Hex()) {
						response.JSON(rw, http.StatusForbidden, nil, "403020")
					}
				}
			}

			if len(contentDetail.LimitUserTag) > 0 {
				foundTag := false
				for _, tag := range contentDetail.LimitUserTag {
					if foundTag = request.CheckUserTag(r, tag); foundTag {
						break
					}
				}

				if !foundTag {
					if !help.NewSlice(help.NewSlice(contentDetail.OnlyUserIDNotLimitUserTag).ObjectIDToStrings()).CheckItem(userID.Hex()) {
						response.JSON(rw, http.StatusForbidden, nil, "403030")
					}
				}
			}

			if !contentDetail.StartTime.IsZero() && !contentDetail.EndTime.IsZero() && !(contentDetail.StartTime.Before(now) && contentDetail.EndTime.After(now)) {
				response.JSON(rw, http.StatusForbidden, nil, "403040")
			} else if !contentDetail.StartTime.IsZero() && !contentDetail.StartTime.Before(now) {
				response.JSON(rw, http.StatusForbidden, nil, "403041")
			} else if !contentDetail.EndTime.IsZero() && !contentDetail.EndTime.After(now) {
				response.JSON(rw, http.StatusForbidden, nil, "403042")
			}

			if !contentDetail.HideDetail ||
				(contentDetail.HideDetail &&
					len(contentDetail.OnlyUserIDShowDetail) > 0 &&
					help.NewSlice(help.NewSlice(contentDetail.OnlyUserIDShowDetail).ObjectIDToStrings()).CheckItem(userID.Hex())) {
				response.JSON(rw, 0, contentDetail, "")
			} else {
				response.JSON(rw, http.StatusForbidden, nil, "403050")
			}
		} else {
			log.Error(err)
			response.JSON(rw, http.StatusServiceUnavailable, nil, "")
		}
	} else {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	}
}

func (self *Service) Route_GET_AddOnlyUseridShowDetail(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var requestOIDR predefined.RequestServeObjectIDRequired
	if err := validator.FormStruct(&requestOIDR, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	contentModel := content.NewModel(self.M)
	if ur, err := contentModel.AddOnlyUserIDShowDetail(r.Context(), requestOIDR.ObjectID, claims.UserID); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}
