package follow

import (
	"net/http"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/plugins/content/database/model/follow_relationship"
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

func (self *Service) Route_POST_Relationship(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestSFRC predefined.RequestServeFollowRelationshipCreate
	if err := validator.FormStruct(&requestSFRC, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	relModel := follow_relationship.NewModel(self.M)

	if _, err := relModel.CreateRelationship(r.Context(), requestSFRC.Type, claims.UserID, requestSFRC.ObjectID, requestSFRC.ObjectOwnerUserID, requestSFRC.Stealth, requestSFRC.ExtensionID); err != nil {
		log.Error(err)
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_DELETE_Relationship(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var requestSFRD predefined.RequestServeFollowRelationshipDelete
	if err := validator.FormStruct(&requestSFRD, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	relModel := follow_relationship.NewModel(self.M)

	if err := relModel.DeleteRelationship(r.Context(), requestSFRD.Type, claims.UserID, requestSFRD.ObjectID); err != nil {
		log.Error(err)
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_DELETE_RelationshipById(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var requestSID predefined.RequestServeIDRequired
	if err := validator.FormStruct(&requestSID, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	relModel := follow_relationship.NewModel(self.M)

	if err := relModel.DeleteRelationshipByID(r.Context(), requestSID.ID); err != nil {
		log.Error(err)
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_GET_Relationships(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var requestSFRL predefined.RequestServeFollowRelationshipList
	if err := validator.FormStruct(&requestSFRL, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	relModel := follow_relationship.NewModel(self.M)
	filter := relModel.FilterByUserWithType(requestSFRL.UserID, requestSFRL.Type)

	cnt, _ := relModel.CountDocuments(r.Context(), filter)
	pg := pagination.Parse(r, cnt)

	projection := bson.D{
		{"_id", 1},
		{"type", 1},
		{"user_id", 1},
		{"object_id", 1},
		{"object_owner_user_id", 1},
		{"extension_id", 1},
		{"mutual", 1},
		{"create_time", 1},
	}

	opt := options.Find().SetSort(bson.D{
		{"create_time", -1},
	}).SetProjection(projection).SetSkip(pg.SkipNum).SetLimit(pg.PageSize)

	if cur, err := relModel.Find(r.Context(), filter, opt); err != nil {
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

func (self *Service) Route_GET_BeRelationships(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var requestSFRL predefined.RequestServeFollowBeRelationshipList
	if err := validator.FormStruct(&requestSFRL, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	relModel := follow_relationship.NewModel(self.M)
	filter := relModel.FilterByObjectIDWithType(requestSFRL.ObjectID, requestSFRL.Type)

	cnt, _ := relModel.CountDocuments(r.Context(), filter)
	pg := pagination.Parse(r, cnt)

	projection := bson.D{
		{"_id", 1},
		{"type", 1},
		{"user_id", 1},
		{"object_id", 1},
		{"object_owner_user_id", 1},
		{"extension_id", 1},
		{"mutual", 1},
		{"create_time", 1},
	}

	opt := options.Find().SetSort(bson.D{
		{"create_time", -1},
	}).SetProjection(projection).SetSkip(pg.SkipNum).SetLimit(pg.PageSize)

	if cur, err := relModel.Find(r.Context(), filter, opt); err != nil {
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

func (self *Service) Route_GET_RelationshipTotal(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var requestSFRT predefined.RequestServeFollowRelationshipTotal
	if err := validator.FormStruct(&requestSFRT, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	relModel := follow_relationship.NewModel(self.M)
	filter := relModel.FilterByUserWithType(requestSFRT.UserID, requestSFRT.Type)

	if requestSFRT.ObjectOwnerUserID != primitive.NilObjectID {
		filter = append(filter, relModel.FilterByObjectOwnerUserID(requestSFRT.ObjectOwnerUserID)...)
	}

	if requestSFRT.ExtensionID != primitive.NilObjectID {
		filter = append(filter, relModel.FilterByExtensionID(requestSFRT.ExtensionID)...)
	}

	if cnt, err := relModel.CountDocuments(r.Context(), filter); err != nil {
		log.Error(err)
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, cnt, "")
	}
}

func (self *Service) Route_GET_BeRelationshipTotal(rw http.ResponseWriter, r *http.Request) {
	claims := request.GetClaims(r)
	if claims == nil {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var requestSFBRT predefined.RequestServeFollowBeRelationshipTotal
	if err := validator.FormStruct(&requestSFBRT, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	relModel := follow_relationship.NewModel(self.M)
	filter := relModel.FilterByObjectIDWithType(requestSFBRT.ObjectID, requestSFBRT.Type)

	if requestSFBRT.ObjectOwnerUserID != primitive.NilObjectID {
		filter = append(filter, relModel.FilterByObjectOwnerUserID(requestSFBRT.ObjectOwnerUserID)...)
	}

	if requestSFBRT.ExtensionID != primitive.NilObjectID {
		filter = append(filter, relModel.FilterByExtensionID(requestSFBRT.ExtensionID)...)
	}

	if cnt, err := relModel.CountDocuments(r.Context(), filter); err != nil {
		log.Error(err)
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, cnt, "")
	}
}
