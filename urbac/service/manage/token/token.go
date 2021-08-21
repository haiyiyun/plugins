package token

import (
	"context"
	"net/http"

	"github.com/haiyiyun/plugins/urbac/database/model"
	"github.com/haiyiyun/plugins/urbac/database/model/token"
	"github.com/haiyiyun/utils/http/pagination"
	"github.com/haiyiyun/utils/http/request"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/validator"
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
		cur.All(context.TODO(), &items)
	}

	rpr := response.ResponsePaginationResult{
		Total: cnt,
		Items: items,
	}

	response.JSON(rw, 0, rpr, "")
}

func (self *Service) Route_DELETE_Delete(rw http.ResponseWriter, r *http.Request) {
	vs, _ := request.ParseDeleteForm(r)
	tokenIDHex := vs.Get("_id")
	tokenString := vs.Get("token")

	valid := &validator.Validation{}
	valid.Required(tokenIDHex).Key("_id").Message("_id字段为空")
	valid.Required(tokenString).Key("token").Message("token字段为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	if tokenID, err := primitive.ObjectIDFromHex(tokenIDHex); err == nil {
		tokenModel := token.NewModel(self.M)
		tokenModel.DeleteOne(context.TODO(), tokenModel.FilterByID(tokenID))
		cacheKey := "claims.valid." + tokenString
		self.Cache.Delete(cacheKey)
		response.JSON(rw, 0, nil, "")

		return
	}

	response.JSON(rw, http.StatusBadRequest, nil, "")
}
