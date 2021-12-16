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

	contentModel := content.NewModel(self.M)
	filterContent := contentModel.FilterNormalContent()

	discussModel := discuss.NewModel(self.M)

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
		{"limit_all_discuss_num", 1},
		{"limit_publish_user_discuss_num", 1},
		{"limit_user_discuss_num", 1},
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
			//是否禁止评论
			if cont.ForbidDiscuss {
				response.JSON(rw, http.StatusForbidden, nil, "403010")
				return
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
						response.JSON(rw, http.StatusForbidden, nil, "403020")
						return
					}
				}
			}

			//限制发布者评论数量处理
			if cont.LimitPublishUserDiscussNum == -1 {
				if cont.PublishUserID == userID {
					response.JSON(rw, http.StatusForbidden, nil, "403030")
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
						response.JSON(rw, http.StatusForbidden, nil, "403031")
						return
					}
				}
			}

			//限制非发布者评论数量处理
			if cont.LimitUserDiscussNum == -1 {
				if cont.PublishUserID != userID {
					response.JSON(rw, http.StatusForbidden, nil, "403040")
					return
				}
			} else if cont.LimitUserDiscussNum > 0 {
				filterLUD := bson.D{
					{"publish_user_id", bson.D{
						{"$ne", cont.PublishUserID},
					}},
				}
				if cntLUD, err := discussModel.CountDocuments(r.Context(), filterLUD); err != nil && err != mongo.ErrNoDocuments {
					log.Error(err)
					response.JSON(rw, http.StatusBadRequest, nil, "400030")
					return
				} else {
					if cntLUD >= int64(cont.LimitUserDiscussNum) {
						response.JSON(rw, http.StatusForbidden, nil, "403041")
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
