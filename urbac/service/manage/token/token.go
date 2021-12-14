package token

import (
	"context"
	"net/http"

	"github.com/haiyiyun/plugins/urbac/database/model"
	"github.com/haiyiyun/plugins/urbac/database/model/token"
	"github.com/haiyiyun/plugins/urbac/predefined"
	"github.com/haiyiyun/utils/http/pagination"
	"github.com/haiyiyun/utils/http/request"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/utils/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Show token list
func (self *Service) Route_GET_Index(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tokenModel := token.NewModel(self.M)
	filter := bson.D{}

	tokenTypes := r.Form["token_type[]"]
	if len(tokenTypes) > 0 {
		filter = append(filter, bson.E{
			"token_type", bson.D{
				{"$in", tokenTypes},
			},
		})
	}

	if userIDHex := r.FormValue("user_id"); userIDHex != "" {
		if userID, err := primitive.ObjectIDFromHex(userIDHex); err == nil {
			filter = append(filter, bson.E{"user_id", userID})
		}
	}

	if userName := r.FormValue("user_name"); userName != "" {
		filter = append(filter, bson.E{"user_name", userName})
	}

	if ip := r.FormValue("ip"); ip != "" {
		filter = append(filter, bson.E{"sign_info.ip", ip})
	}

	cnt, _ := tokenModel.CountDocuments(context.Background(), filter)
	pg := pagination.Parse(r, cnt)

	opt := options.Find().SetSort(bson.D{
		{"create_time", -1},
	}).SetProjection(bson.D{}).SetSkip(pg.SkipNum).SetLimit(pg.PageSize)

	items := []model.Token{}
	if cur, err := tokenModel.Find(context.Background(), filter, opt); err == nil {
		cur.All(r.Context(), &items)
	}

	rpr := response.ResponsePaginationResult{
		Total: cnt,
		Items: items,
	}

	response.JSON(rw, 0, rpr, "")
}

func (self *Service) Route_DELETE_Delete(rw http.ResponseWriter, r *http.Request) {
	vs, _ := request.ParseDeleteForm(r)

	var requestMTD predefined.RequestManageTokenDelete
	if err := validator.FormStruct(&requestMTD, vs); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	tokenModel := token.NewModel(self.M)
	if dr, err := tokenModel.DeleteOne(r.Context(), tokenModel.FilterByID(requestMTD.ObjectID)); err == nil && dr.DeletedCount > 0 {
		cacheKey := "claims.valid." + requestMTD.ObjectID.Hex()
		self.Cache.Delete(cacheKey)
		response.JSON(rw, 0, nil, "")
	} else {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	}
}
