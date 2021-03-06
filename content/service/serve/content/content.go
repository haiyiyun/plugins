package content

import (
	"context"
	"net/http"
	"time"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/mongodb"
	"github.com/haiyiyun/mongodb/geometry"
	"github.com/haiyiyun/plugins/content/database/model"
	"github.com/haiyiyun/plugins/content/database/model/category"
	"github.com/haiyiyun/plugins/content/database/model/content"
	"github.com/haiyiyun/plugins/content/database/model/follow_content"
	"github.com/haiyiyun/plugins/content/database/model/follow_relationship"
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
					if cont.LimitAssociateType != requestCC.AssociateType {
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
					if cnt, err := contentModel.CountDocuments(r.Context(), filterLimit); err != nil {
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
							return
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
							return
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
							return
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
							return
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
							return
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
		//关注处理
		go func(mgo mongodb.Mongoer, userID, subjectID, contentID primitive.ObjectID, publishType, associateType int, associateID primitive.ObjectID) {
			frModel := follow_relationship.NewModel(mgo)

			//关注人处理
			filter := frModel.FilterByObjectID(userID)
			filter = append(filter, frModel.FilterByType(predefined.FollowTypeUser)...)
			ctx := context.Background()
			if cur, err := frModel.Find(ctx, filter, options.Find().SetProjection(bson.D{
				{"_id", 1},
				{"extension_id", 1},
				{"user_id", 1},
			})); err != nil {
				log.Error(err)
			} else {
				var frls []model.FollowRelationship
				if err := cur.All(ctx, &frls); err != nil {
					log.Error(err)
				} else {
					for _, fr := range frls {
						go func(mgo mongodb.Mongoer, fr model.FollowRelationship, contentID primitive.ObjectID, typ int) {
							fcModel := follow_content.NewModel(mgo)
							//判断是否已经创建
							filter := fcModel.FilterByFollowRelationshipID(fr.ID)
							filter = append(filter, fcModel.FilterByUserID(fr.UserID)...)
							filter = append(filter, fcModel.FilterByContentID(contentID)...)
							if cnt, err := fcModel.CountDocuments(context.Background(), filter); err != nil {
								log.Error(err)
							} else {
								if cnt == 0 {
									if _, err := fcModel.Create(context.Background(), &model.FollowContent{
										FollowRelationshipID: fr.ID,
										Type:                 typ,
										UserID:               fr.UserID,
										ContentID:            contentID,
										ExtensionID:          fr.ExtensionID,
									}); err != nil {
										log.Error(err)
									}
								}
							}
						}(mgo, fr, contentID, publishType)
					}
				}
			}

			//关注主题处理
			if subjectID != primitive.NilObjectID {
				filterSID := frModel.FilterByObjectID(subjectID)
				filterSID = append(filterSID, frModel.FilterByType(predefined.FollowTypeSubject)...)
				ctxSID := context.Background()
				if cur, err := frModel.Find(ctxSID, filterSID, options.Find().SetProjection(bson.D{
					{"_id", 1},
					{"extension_id", 1},
					{"user_id", 1},
				})); err != nil {
					log.Error(err)
				} else {
					var frls []model.FollowRelationship
					if err := cur.All(ctxSID, &frls); err != nil {
						log.Error(err)
					} else {
						for _, fr := range frls {
							go func(mgo mongodb.Mongoer, fr model.FollowRelationship, contentID primitive.ObjectID, typ int) {
								fcModel := follow_content.NewModel(mgo)
								//判断是否已经创建
								filter := fcModel.FilterByFollowRelationshipID(fr.ID)
								filter = append(filter, fcModel.FilterByUserID(fr.UserID)...)
								filter = append(filter, fcModel.FilterByContentID(contentID)...)
								if cnt, err := fcModel.CountDocuments(context.Background(), filter); err != nil {
									log.Error(err)
								} else {
									if cnt == 0 {
										if _, err := fcModel.Create(context.Background(), &model.FollowContent{
											FollowRelationshipID: fr.ID,
											Type:                 typ,
											UserID:               fr.UserID,
											ContentID:            contentID,
											ExtensionID:          fr.ExtensionID,
										}); err != nil {
											log.Error(err)
										}
									}
								}
							}(mgo, fr, contentID, publishType)
						}
					}
				}
			}

			//关注关联处理
			if associateID != primitive.NilObjectID {
				//关联类型过滤
				if associateType != predefined.ContentAssociateTypeSelf && associateType != predefined.ContentAssociateTypeForward {
					filterAID := frModel.FilterByObjectID(associateID)
					filterAID = append(filterAID, frModel.FilterByType(publishType)...)
					ctxAID := context.Background()
					if cur, err := frModel.Find(ctxAID, filterAID, options.Find().SetProjection(bson.D{
						{"_id", 1},
						{"extension_id", 1},
						{"user_id", 1},
					})); err != nil {
						log.Error(err)
					} else {
						var frls []model.FollowRelationship
						if err := cur.All(ctxAID, &frls); err != nil {
							log.Error(err)
						} else {
							for _, fr := range frls {
								go func(mgo mongodb.Mongoer, fr model.FollowRelationship, contentID primitive.ObjectID, typ int) {
									fcModel := follow_content.NewModel(mgo)
									//判断是否已经创建
									filter := fcModel.FilterByFollowRelationshipID(fr.ID)
									filter = append(filter, fcModel.FilterByUserID(fr.UserID)...)
									filter = append(filter, fcModel.FilterByContentID(contentID)...)
									if cnt, err := fcModel.CountDocuments(context.Background(), filter); err != nil {
										log.Error(err)
									} else {
										if cnt == 0 {
											if _, err := fcModel.Create(context.Background(), &model.FollowContent{
												FollowRelationshipID: fr.ID,
												Type:                 typ,
												UserID:               fr.UserID,
												ContentID:            contentID,
												ExtensionID:          fr.ExtensionID,
											}); err != nil {
												log.Error(err)
											}
										}
									}
								}(mgo, fr, contentID, publishType)
							}
						}
					}
				}
			}
		}(self.M, userID, ctnt.SubjectID, ior.InsertedID.(primitive.ObjectID), ctnt.PublishType, ctnt.AssociateType, ctnt.AssociateID)

		response.JSON(rw, 0, ior.InsertedID, "")
	}
}

func (self *Service) Route_DELETE_Content(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var requestSCD predefined.RequestServeContentDelete
	if err := validator.FormStruct(&requestSCD, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	contentModel := content.NewModel(self.M)
	filter := contentModel.FilterByPublishUserID(claims.UserID)
	filter = append(filter, contentModel.FilterByID(requestSCD.ID)...)

	if ur, err := contentModel.Set(r.Context(), filter, bson.D{
		{"status", predefined.PublishStatusDelete},
	}); err != nil || ur.ModifiedCount == 0 {
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

	var requestCL predefined.RequestServeContentList
	if err := validator.FormStruct(&requestCL, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	contentModel := content.NewModel(self.M)
	filter := contentModel.FilterNormalContent()

	project := bson.D{
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
		{"in_readed_user", bson.D{
			{"$in", bson.A{claims.UserID, `$readed_user`}},
		}},
		{"readed_user_total", bson.D{
			{"$size", `$readed_user`},
		}},
		{"in_wanted_user", bson.D{
			{"$in", bson.A{claims.UserID, `$wanted_user`}},
		}},
		{"wanted_user_total", bson.D{
			{"$size", `$wanted_user`},
		}},
		{"in_liked_user", bson.D{
			{"$in", bson.A{claims.UserID, `$liked_user`}},
		}},
		{"liked_user_total", bson.D{
			{"$size", `$liked_user`},
		}},
		{"in_hated_user", bson.D{
			{"$in", bson.A{claims.UserID, `$hated_user`}},
		}},
		{"hated_user_total", bson.D{
			{"$size", `$hated_user`},
		}},
		{"guise", 1},
		{"anti_guise_user", 1},
		{"start_time", 1},
		{"end_time", 1},
		{"extra_data", 1},
		{"status", 1},
		{"discuss_estimate_total", 1},
		{"discuss_estimate_evaluation_total", 1},
		{"create_time", 1},
		{"update_time", 1},
	}

	var limit int64
	var skip int64
	sort := bson.D{
		{"create_time", -1},
	}

	var cnt int64

	if requestCL.ID != primitive.NilObjectID {
		filter = append(filter, contentModel.FilterByID(requestCL.ID)...)
		limit = 1
		cnt = 1
	} else {
		coordinates := geometry.PointCoordinates{
			requestCL.Longitude, requestCL.Latitude,
		}

		if len(requestCL.Types) > 0 {
			filter = append(filter, contentModel.FilterByTypes(requestCL.Types)...)
		}

		if requestCL.PublishType > 0 {
			filter = append(filter, contentModel.FilterByPublishType(requestCL.PublishType)...)
		}

		if requestCL.Visibility == predefined.VisibilityTypeSelf {
			filter = append(filter, contentModel.FilterByPublishUserID(claims.UserID)...)
		} else {
			if len(requestCL.PublishUserID) > 0 {
				filter = append(filter, contentModel.FilterByPublishUserIDs(requestCL.PublishUserID)...)
			}

			if requestCL.Visibility != predefined.VisibilityTypeAll {
				filter = append(filter, contentModel.FilterByVisibilityOrAll(requestCL.Visibility)...)
			} else {
				filter = append(filter, contentModel.FilterByVisibility(requestCL.Visibility)...)
			}
		}

		if len(requestCL.ExcludePublishUserID) > 0 {
			filter = append(filter, contentModel.FilterByExcludePublishUserIDs(requestCL.ExcludePublishUserID)...)
		}

		if requestCL.AssociateType > 0 {
			filter = append(filter, contentModel.FilterByAssociateType(requestCL.AssociateType)...)
		}

		if requestCL.AssociateID != primitive.NilObjectID {
			//判断关联是否有限制
			filterAssociateContent := contentModel.FilterNormalContent()
			filterAssociateContent = append(filterAssociateContent, contentModel.FilterByID(requestCL.AssociateID)...)
			associateContentSR := contentModel.FindOne(r.Context(), filterAssociateContent, options.FindOne().SetProjection(bson.D{
				{"hide_discuss", 1},
				{"only_user_id_show_discuss", 1},
			}))

			if associateContentSR.Err() != nil {
				log.Error(associateContentSR.Err())
				response.JSON(rw, http.StatusServiceUnavailable, nil, "503000")
				return
			}

			var associateContent model.Content
			if err := associateContentSR.Decode(&associateContent); err != nil {
				log.Error(err)
				response.JSON(rw, http.StatusServiceUnavailable, nil, "503001")
				return
			} else {
				if associateContent.HideDetail {
					if len(associateContent.OnlyUserIDShowDetail) == 0 || (len(associateContent.OnlyUserIDShowDetail) > 0 && !help.NewSlice(help.NewSlice(associateContent.OnlyUserIDShowDetail).ObjectIDToStrings()).CheckItem(claims.UserID.Hex())) {
						response.JSON(rw, http.StatusForbidden, nil, "403000")
						return
					}
				}
			}

			filter = append(filter, contentModel.FilterByAssociateID(requestCL.AssociateID)...)
		}

		if requestCL.EmptyCategoryID {
			filter = append(filter, contentModel.FilterByCategoryID(primitive.NilObjectID)...)
		} else {
			if requestCL.CategoryID != primitive.NilObjectID {
				filter = append(filter, contentModel.FilterByCategoryID(requestCL.CategoryID)...)
			}
		}

		if requestCL.EmptySubjectID {
			filter = append(filter, contentModel.FilterBySubjectID(primitive.NilObjectID)...)
		} else {
			if requestCL.SubjectID != primitive.NilObjectID {
				filter = append(filter, contentModel.FilterBySubjectID(requestCL.SubjectID)...)
			}
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
			filter = append(filter, contentModel.FilterByStartTimeLte(requestCL.StartTime.Time)...)
		}

		if !requestCL.EndTime.Time.IsZero() {
			filter = append(filter, contentModel.FilterByEndTimeGte(requestCL.EndTime.Time)...)
		}

		if requestCL.InTime {
			now := time.Now()
			filter = append(filter, contentModel.FilterByStartTimeLte(now)...)
			filter = append(filter, contentModel.FilterByEndTimeGte(now)...)
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

		cnt, _ = contentModel.CountDocuments(r.Context(), filter)
		pg := pagination.Parse(r, cnt)
		limit = pg.PageSize
		skip = pg.SkipNum
	}

	if cur, err := contentModel.Aggregate(r.Context(), mongo.Pipeline{
		{{"$match", filter}},
		{{"$project", project}},
		{{"$sort", sort}},
		{{"$skip", skip}},
		{{"$limit", limit}},
	}); err != nil {
		log.Error(err)
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
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

	var requestCD predefined.RequestServeContentDetail
	if err := validator.FormStruct(&requestCD, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	userID := claims.UserID
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
						return
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
						return
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
						return
					}
				}
			}

			if !contentDetail.StartTime.IsZero() && !contentDetail.EndTime.IsZero() && !(contentDetail.StartTime.Before(now) && contentDetail.EndTime.After(now)) {
				response.JSON(rw, http.StatusForbidden, nil, "403040")
				return
			} else if !contentDetail.StartTime.IsZero() && !contentDetail.StartTime.Before(now) {
				response.JSON(rw, http.StatusForbidden, nil, "403041")
				return
			} else if !contentDetail.EndTime.IsZero() && !contentDetail.EndTime.After(now) {
				response.JSON(rw, http.StatusForbidden, nil, "403042")
				return
			}

			if !contentDetail.HideDetail ||
				(contentDetail.HideDetail &&
					len(contentDetail.OnlyUserIDShowDetail) > 0 &&
					help.NewSlice(help.NewSlice(contentDetail.OnlyUserIDShowDetail).ObjectIDToStrings()).CheckItem(userID.Hex())) {
				response.JSON(rw, 0, contentDetail, "")

				//判断是否已经添加过已读
				if !help.NewSlice(help.NewSlice(contentDetail.ReadedUser).ObjectIDToStrings()).CheckItem(userID.Hex()) {
					//添加已读user_id
					go contentModel.AddReadedUser(context.Background(), requestCD.ID, userID)
					//批量更新所有关注内容的已读时间
					go func(mgo mongodb.Mongoer, userID, contentID primitive.ObjectID) {
						fcModel := follow_content.NewModel(mgo)
						fcModel.SetAllReadedTimeByUserIDAndContentID(context.Background(), userID, contentID, time.Now())
					}(self.M, userID, requestCD.ID)
				}
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

func (self *Service) Route_POST_ReadedUser(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestOIDR predefined.RequestServeObjectIDRequired
	if err := validator.FormStruct(&requestOIDR, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	contentModel := content.NewModel(self.M)
	if ur, err := contentModel.AddReadedUser(r.Context(), requestOIDR.ObjectID, claims.UserID); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_DELETE_ReadedUser(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestOIDR predefined.RequestServeObjectIDRequired
	if err := validator.FormStruct(&requestOIDR, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	contentModel := content.NewModel(self.M)
	if ur, err := contentModel.DeleteReadedUser(r.Context(), requestOIDR.ObjectID, claims.UserID); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_POST_WantedUser(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestOIDR predefined.RequestServeObjectIDRequired
	if err := validator.FormStruct(&requestOIDR, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	contentModel := content.NewModel(self.M)
	if ur, err := contentModel.AddWantedUser(r.Context(), requestOIDR.ObjectID, claims.UserID); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_DELETE_WantedUser(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestOIDR predefined.RequestServeObjectIDRequired
	if err := validator.FormStruct(&requestOIDR, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	contentModel := content.NewModel(self.M)
	if ur, err := contentModel.DeleteWantedUser(r.Context(), requestOIDR.ObjectID, claims.UserID); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_POST_LikedUser(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestOIDR predefined.RequestServeObjectIDRequired
	if err := validator.FormStruct(&requestOIDR, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	contentModel := content.NewModel(self.M)
	if ur, err := contentModel.AddLikedUser(r.Context(), requestOIDR.ObjectID, claims.UserID); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_DELETE_LikedUser(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestOIDR predefined.RequestServeObjectIDRequired
	if err := validator.FormStruct(&requestOIDR, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	contentModel := content.NewModel(self.M)
	if ur, err := contentModel.DeleteLikedUser(r.Context(), requestOIDR.ObjectID, claims.UserID); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_POST_HatedUser(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestOIDR predefined.RequestServeObjectIDRequired
	if err := validator.FormStruct(&requestOIDR, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	contentModel := content.NewModel(self.M)
	if ur, err := contentModel.AddHatedUser(r.Context(), requestOIDR.ObjectID, claims.UserID); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_DELETE_HatedUser(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestOIDR predefined.RequestServeObjectIDRequired
	if err := validator.FormStruct(&requestOIDR, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	contentModel := content.NewModel(self.M)
	if ur, err := contentModel.DeleteHatedUser(r.Context(), requestOIDR.ObjectID, claims.UserID); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_POST_OnlyUseridShowDetail(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestOIDR predefined.RequestServeObjectIDRequired
	if err := validator.FormStruct(&requestOIDR, r.Form); err != nil {
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

func (self *Service) Route_POST_Description(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestSCU predefined.RequestServeContentUpdateDescription
	if err := validator.FormStruct(&requestSCU, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	contentModel := content.NewModel(self.M)
	filter := contentModel.FilterByPublishUserID(claims.UserID)
	filter = append(filter, contentModel.FilterByID(requestSCU.ObjectID)...)
	if ur, err := contentModel.Set(r.Context(), filter, bson.D{
		{"description", requestSCU.Description},
	}); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_POST_Visibility(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestSCU predefined.RequestServeContentUpdateVisibility
	if err := validator.FormStruct(&requestSCU, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	contentModel := content.NewModel(self.M)
	filter := contentModel.FilterByPublishUserID(claims.UserID)
	filter = append(filter, contentModel.FilterByID(requestSCU.ObjectID)...)
	if ur, err := contentModel.Set(r.Context(), filter, bson.D{
		{"visibility", requestSCU.Visibility},
	}); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_POST_ForbidForward(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestSCU predefined.RequestServeContentUpdateForbidForward
	if err := validator.FormStruct(&requestSCU, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	contentModel := content.NewModel(self.M)
	filter := contentModel.FilterByPublishUserID(claims.UserID)
	filter = append(filter, contentModel.FilterByID(requestSCU.ObjectID)...)
	if ur, err := contentModel.Set(r.Context(), filter, bson.D{
		{"forbid_forward", requestSCU.ForbidForward},
	}); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_POST_ForbidDownload(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestSCU predefined.RequestServeContentUpdateForbidDownload
	if err := validator.FormStruct(&requestSCU, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	contentModel := content.NewModel(self.M)
	filter := contentModel.FilterByPublishUserID(claims.UserID)
	filter = append(filter, contentModel.FilterByID(requestSCU.ObjectID)...)
	if ur, err := contentModel.Set(r.Context(), filter, bson.D{
		{"forbid_download", requestSCU.ForbidDownload},
	}); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_POST_ForbidDiscuss(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestSCU predefined.RequestServeContentUpdateForbidDiscuss
	if err := validator.FormStruct(&requestSCU, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	contentModel := content.NewModel(self.M)
	filter := contentModel.FilterByPublishUserID(claims.UserID)
	filter = append(filter, contentModel.FilterByID(requestSCU.ObjectID)...)
	if ur, err := contentModel.Set(r.Context(), filter, bson.D{
		{"forbid_discuss", requestSCU.ForbidDiscuss},
	}); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_POST_StartTime(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestSCU predefined.RequestServeContentUpdateStartTime
	if err := validator.FormStruct(&requestSCU, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	contentModel := content.NewModel(self.M)
	filter := contentModel.FilterByPublishUserID(claims.UserID)
	filter = append(filter, contentModel.FilterByID(requestSCU.ObjectID)...)
	if ur, err := contentModel.Set(r.Context(), filter, bson.D{
		{"start_time", requestSCU.StartTime.Time},
	}); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_POST_EndTime(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestSCU predefined.RequestServeContentUpdateEndTime
	if err := validator.FormStruct(&requestSCU, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	contentModel := content.NewModel(self.M)
	filter := contentModel.FilterByPublishUserID(claims.UserID)
	filter = append(filter, contentModel.FilterByID(requestSCU.ObjectID)...)
	if ur, err := contentModel.Set(r.Context(), filter, bson.D{
		{"end_time", requestSCU.EndTime.Time},
	}); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_POST_ExtraData(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestSCU predefined.RequestServeContentUpdateExtraData
	if err := validator.FormStruct(&requestSCU, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	contentModel := content.NewModel(self.M)
	filter := contentModel.FilterByPublishUserID(claims.UserID)
	filter = append(filter, contentModel.FilterByID(requestSCU.ObjectID)...)
	if ur, err := contentModel.Set(r.Context(), filter, bson.D{
		{"extra_data", requestSCU.ExtraData},
	}); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}
