package discuss

import (
	"context"
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

	discussModel := discuss.NewModel(self.M)

	//判断被回复的评论是否存在
	if requestDC.ReplyDiscussID != primitive.NilObjectID {
		filter := discussModel.FilterNormalDiscuss()
		filter = append(filter, discussModel.FilterByID(requestDC.ReplyDiscussID)...)
		if cnt, err := discussModel.CountDocuments(r.Context(), filter); cnt == 0 {
			log.Error(err)
			response.JSON(rw, http.StatusBadRequest, nil, "400404")
			return
		}
	}

	contentModel := content.NewModel(self.M)
	filterContent := contentModel.FilterNormalContent()

	//判断对应type的object_id是否存在,并处理相关限制
	switch requestDC.Type {
	case predefined.DiscussTypeDynamic:
		contentType := predefined.ContentPublishTypeDynamic
		filterContent = append(filterContent, contentModel.FilterByID(requestDC.ObjectID)...)
		filterContent = append(filterContent, contentModel.FilterByPublishType(contentType)...)
	case predefined.DiscussTypeArticle:
		contentType := predefined.ContentPublishTypeArticle
		filterContent = append(filterContent, contentModel.FilterByID(requestDC.ObjectID)...)
		filterContent = append(filterContent, contentModel.FilterByPublishType(contentType)...)
	case predefined.DiscussTypeQuestion:
		contentType := predefined.ContentPublishTypeQuestion
		filterContent = append(filterContent, contentModel.FilterByID(requestDC.ObjectID)...)
		filterContent = append(filterContent, contentModel.FilterByPublishType(contentType)...)
	case predefined.DiscussTypeAnswer:
		contentType := predefined.ContentPublishTypeAnswer
		filterContent = append(filterContent, contentModel.FilterByID(requestDC.ObjectID)...)
		filterContent = append(filterContent, contentModel.FilterByPublishType(contentType)...)
	}

	if sr := contentModel.FindOne(r.Context(), filterContent, options.FindOne().SetProjection(bson.D{
		{"publish_user_id", 1},
		{"forbid_discuss", 1},
		{"only_user_id_discuss", 1},
		{"only_publish_user_id_can_reply_discuss", 1},
		{"only_publish_user_id_can_not_reply_discuss", 1},
		{"limit_all_discuss_num", 1},
		{"limit_publish_user_discuss_num", 1},
		{"limit_not_publish_user_all_user_discuss_num", 1},
		{"limit_not_publish_user_every_user_discuss_num", 1},
	})); sr.Err() != nil {
		log.Error(sr.Err())
		response.JSON(rw, http.StatusBadRequest, nil, "400414")
		return
	} else {
		var cont model.Content
		if err := sr.Decode(&cont); err != nil {
			log.Error(err)
			response.JSON(rw, http.StatusBadRequest, nil, "400000")
			return
		} else {
			//是否禁止评论
			if cont.ForbidDiscuss {
				response.JSON(rw, http.StatusForbidden, nil, "403010")
				return
			}

			//判断是否限制user_id可以评论
			if cont.OnlyUserIDDiscuss != nil && len(cont.OnlyUserIDDiscuss) > 0 {
				onlyUserIDDiscussHexs := []string{}
				for _, uid := range cont.OnlyUserIDDiscuss {
					onlyUserIDDiscussHexs = append(onlyUserIDDiscussHexs, uid.Hex())
				}

				if !help.NewSlice(onlyUserIDDiscussHexs).CheckItem(userID.Hex()) {
					response.JSON(rw, http.StatusForbidden, nil, "403020")
					return
				}
			}

			if requestDC.ReplyDiscussID != primitive.NilObjectID {
				//判断是否限制user_id可以回复评论
				if cont.OnlyUserIDCanReplyDiscuss != nil && len(cont.OnlyUserIDCanReplyDiscuss) > 0 {
					onlyUserIDCanReplyDiscuss := []string{}
					for _, uid := range cont.OnlyUserIDCanReplyDiscuss {
						onlyUserIDCanReplyDiscuss = append(onlyUserIDCanReplyDiscuss, uid.Hex())
					}

					if !help.NewSlice(onlyUserIDCanReplyDiscuss).CheckItem(userID.Hex()) {
						response.JSON(rw, http.StatusForbidden, nil, "403030")
						return
					}
				}

				//判断是否限制user_id不可以回复评论
				if cont.OnlyUserIDCanNotReplyDiscuss != nil && len(cont.OnlyUserIDCanNotReplyDiscuss) > 0 {
					onlyUserIDCanNotReplyDiscuss := []string{}
					for _, uid := range cont.OnlyUserIDCanNotReplyDiscuss {
						onlyUserIDCanNotReplyDiscuss = append(onlyUserIDCanNotReplyDiscuss, uid.Hex())
					}

					if help.NewSlice(onlyUserIDCanNotReplyDiscuss).CheckItem(userID.Hex()) {
						response.JSON(rw, http.StatusForbidden, nil, "403040")
						return
					}
				}
			}

			//限制所有评论数处理
			if cont.LimitAllDiscussNum > 0 {
				filterLAD := discussModel.FilterNormalDiscuss()
				filterLAD = append(filterLAD, discussModel.FilterByObjectID(requestDC.ObjectID)...)
				if cntLAD, err := discussModel.CountDocuments(r.Context(), filterLAD); err != nil && err != mongo.ErrNoDocuments {
					log.Error(err)
					response.JSON(rw, http.StatusBadRequest, nil, "400010")
					return
				} else {
					if cntLAD >= int64(cont.LimitAllDiscussNum) {
						response.JSON(rw, http.StatusForbidden, nil, "403050")
						return
					}
				}
			}

			//限制发布者评论数量处理
			if cont.LimitPublishUserDiscussNum == -1 {
				if cont.PublishUserID == userID {
					response.JSON(rw, http.StatusForbidden, nil, "403060")
					return
				}
			} else if cont.LimitPublishUserDiscussNum > 0 {
				filterLPUD := discussModel.FilterByPublishUserID(cont.PublishUserID)
				if cntLPUD, err := discussModel.CountDocuments(r.Context(), filterLPUD); err != nil && err != mongo.ErrNoDocuments {
					log.Error(err)
					response.JSON(rw, http.StatusBadRequest, nil, "400020")
					return
				} else {
					if cntLPUD >= int64(cont.LimitPublishUserDiscussNum) {
						response.JSON(rw, http.StatusForbidden, nil, "403071")
						return
					}
				}
			}

			//限制非发布者评论数量处理
			if cont.LimitNotPublishUserAllUserDiscussNum == -1 {
				if cont.PublishUserID != userID {
					response.JSON(rw, http.StatusForbidden, nil, "403080")
					return
				}
			} else if cont.LimitNotPublishUserAllUserDiscussNum > 0 {
				filterLNPUAUDN := bson.D{
					{"publish_user_id", bson.D{
						{"$ne", cont.PublishUserID},
					}},
				}
				if cntLUD, err := discussModel.CountDocuments(r.Context(), filterLNPUAUDN); err != nil && err != mongo.ErrNoDocuments {
					log.Error(err)
					response.JSON(rw, http.StatusBadRequest, nil, "400030")
					return
				} else {
					if cntLUD >= int64(cont.LimitNotPublishUserAllUserDiscussNum) {
						response.JSON(rw, http.StatusForbidden, nil, "403081")
						return
					}
				}
			} else if cont.LimitNotPublishUserEveryUserDiscussNum > 0 && cont.PublishUserID != userID {
				filterLNPUEUDN := bson.D{
					{"publish_user_id", userID},
				}
				if cntLUD, err := discussModel.CountDocuments(r.Context(), filterLNPUEUDN); err != nil && err != mongo.ErrNoDocuments {
					log.Error(err)
					response.JSON(rw, http.StatusBadRequest, nil, "400031")
					return
				} else {
					if cntLUD >= int64(cont.LimitNotPublishUserEveryUserDiscussNum) {
						response.JSON(rw, http.StatusForbidden, nil, "403082")
						return
					}
				}
			}
		}
	}

	dis := &model.Discuss{
		Type:           requestDC.Type,
		ObjectID:       requestDC.ObjectID,
		PublishUserID:  userID,
		AtUsers:        requestDC.AtUsers,
		ReplyDiscussID: requestDC.ReplyDiscussID,
		Text:           requestDC.Text,
		Location:       geometry.NewPoint(coordinates),
		Visibility:     requestDC.Visibility,
		Evaluation:     requestDC.Evaluation,
		Status:         status,
	}

	if ior, err := discussModel.Create(r.Context(), dis); err != nil || ior.InsertedID == nil {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		go contentModel.Inc(context.Background(), contentModel.FilterByID(requestDC.ObjectID), bson.D{
			{"discuss_estimate_total", 1},
		})

		if requestDC.ReplyDiscussID != primitive.NilObjectID {
			go discussModel.Inc(context.Background(), discussModel.FilterByID(requestDC.ReplyDiscussID), bson.D{
				{"reply_estimate_total", 1},
			})
		}

		response.JSON(rw, 0, ior.InsertedID, "")
	}
}

func (self *Service) Route_DELETE_Discuss(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var requestSDD predefined.RequestServeDiscussDelete
	if err := validator.FormStruct(&requestSDD, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	discussModel := discuss.NewModel(self.M)
	filter := discussModel.FilterByPublishUserID(claims.UserID)
	filter = append(filter, discussModel.FilterByID(requestSDD.ID)...)

	if ur, err := discussModel.Set(r.Context(), filter, bson.D{
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

	discussModel := discuss.NewModel(self.M)
	filter := discussModel.FilterNormalDiscuss()
	filter = append(filter, bson.D{
		{"reply_discuss_id", primitive.NilObjectID},
	}...)

	if requestDL.ID != primitive.NilObjectID {
		filter = append(filter, discussModel.FilterByID(requestDL.ID)...)
	} else {
		if requestDL.ObjectID != primitive.NilObjectID && len(requestDL.Types) > 0 {
			contentModel := content.NewModel(self.M)
			filterContent := contentModel.FilterNormalContent()
			typ := requestDL.Types[0]

			//判断对应type的object_id是否存在,并处理相关限制
			switch typ {
			case predefined.DiscussTypeDynamic:
				contentType := predefined.ContentPublishTypeDynamic
				filterContent = append(filterContent, contentModel.FilterByID(requestDL.ObjectID)...)
				filterContent = append(filterContent, contentModel.FilterByPublishType(contentType)...)
			case predefined.DiscussTypeArticle:
				contentType := predefined.ContentPublishTypeArticle
				filterContent = append(filterContent, contentModel.FilterByID(requestDL.ObjectID)...)
				filterContent = append(filterContent, contentModel.FilterByPublishType(contentType)...)
			case predefined.DiscussTypeQuestion:
				contentType := predefined.ContentPublishTypeQuestion
				filterContent = append(filterContent, contentModel.FilterByID(requestDL.ObjectID)...)
				filterContent = append(filterContent, contentModel.FilterByPublishType(contentType)...)
			case predefined.DiscussTypeAnswer:
				contentType := predefined.ContentPublishTypeAnswer
				filterContent = append(filterContent, contentModel.FilterByID(requestDL.ObjectID)...)
				filterContent = append(filterContent, contentModel.FilterByPublishType(contentType)...)
			}

			if sr := contentModel.FindOne(r.Context(), filterContent, options.FindOne().SetProjection(bson.D{
				{"hide_discuss", 1},
				{"only_user_id_show_discuss", 1},
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
					//是否隐藏评论
					if cont.HideDiscuss {
						//是否有允许查看评论的
						if !(len(cont.OnlyUserIDShowDiscuss) > 0 && help.NewSlice(help.NewSlice(cont.OnlyUserIDShowDiscuss).ObjectIDToStrings()).CheckItem(userID.Hex())) {
							response.JSON(rw, http.StatusForbidden, nil, "403010")
							return
						}
					}
				}
			}

			filter = append(filter, discussModel.FilterByTypes(requestDL.Types)...)
			filter = append(filter, discussModel.FilterByObjectID(requestDL.ObjectID)...)
		}

		if requestDL.Visibility == predefined.VisibilityTypeSelf {
			filter = append(filter, discussModel.FilterByPublishUserID(userID)...)
		} else {
			if len(requestDL.PublishUserID) > 0 {
				filter = append(filter, discussModel.FilterByPublishUserIDs(requestDL.PublishUserID)...)
			}

			if requestDL.Visibility != predefined.VisibilityTypeAll {
				filter = append(filter, discussModel.FilterByVisibilityOrAll(requestDL.Visibility)...)
			} else {
				filter = append(filter, discussModel.FilterByVisibility(requestDL.Visibility)...)
			}
		}

		coordinates := geometry.PointCoordinates{
			requestDL.Longitude, requestDL.Latitude,
		}

		if coordinates != geometry.NilPointCoordinates {
			filter = append(filter, discussModel.FilterByLocation(geometry.NewPoint(coordinates), requestDL.MaxDistance, requestDL.MinDistance)...)
		}
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

func (self *Service) Route_GET_AvgEvaluation(rw http.ResponseWriter, r *http.Request) {
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

	discussModel := discuss.NewModel(self.M)
	match := discussModel.FilterNormalDiscuss()
	match = append(match, discussModel.FilterByObjectID(requestOIDR.ObjectID)...)

	if cur, err := discussModel.Aggregate(r.Context(), mongo.Pipeline{
		{{"$match", match}},
		{{"$group", bson.D{
			{"_id", `$object_id`},
			{"avg_evaluation", bson.D{
				{"$avg", `$evaluation`},
			}},
		}}},
	}); err != nil {
		log.Error(err)
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		var dissAvgs []help.M
		if err := cur.All(r.Context(), &dissAvgs); err != nil {
			log.Error(err)
			response.JSON(rw, http.StatusServiceUnavailable, nil, "")
		} else {
			response.JSON(rw, 0, dissAvgs, "")
		}
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

	discussModel := discuss.NewModel(self.M)
	if ur, err := discussModel.AddLikedUser(r.Context(), requestOIDR.ObjectID, claims.UserID); err != nil || ur.ModifiedCount == 0 {
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

	discussModel := discuss.NewModel(self.M)
	if ur, err := discussModel.DeleteLikedUser(r.Context(), requestOIDR.ObjectID, claims.UserID); err != nil || ur.ModifiedCount == 0 {
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

	discussModel := discuss.NewModel(self.M)
	if ur, err := discussModel.AddHatedUser(r.Context(), requestOIDR.ObjectID, claims.UserID); err != nil || ur.ModifiedCount == 0 {
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

	discussModel := discuss.NewModel(self.M)
	if ur, err := discussModel.DeleteHatedUser(r.Context(), requestOIDR.ObjectID, claims.UserID); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}
