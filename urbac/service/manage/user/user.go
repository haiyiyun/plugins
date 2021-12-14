package user

import (
	"context"
	"net/http"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/plugins/urbac/database/model"
	"github.com/haiyiyun/plugins/urbac/database/model/user"
	"github.com/haiyiyun/plugins/urbac/predefined"

	"github.com/haiyiyun/utils/help"
	"github.com/haiyiyun/utils/http/pagination"
	"github.com/haiyiyun/utils/http/request"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/utils/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (self *Service) Route_GET_Index(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userModel := user.NewModel(self.M)
	filter := bson.D{}

	enables := r.Form["enable[]"]
	if len(enables) > 0 {
		filter = append(filter, bson.E{
			"enable", bson.D{
				{"$in", help.NewSlice(enables).ConvBool()},
			},
		})
	}

	deletes := r.Form["delete[]"]
	if len(deletes) > 0 {
		filter = append(filter, bson.E{
			"delete", bson.D{
				{"$in", help.NewSlice(deletes).ConvBool()},
			},
		})
	}

	if userIDHex := r.FormValue("_id"); userIDHex != "" {
		if userID, err := primitive.ObjectIDFromHex(userIDHex); err == nil {
			filter = append(filter, bson.E{"_id", userID})
		}
	}

	if name := r.FormValue("name"); name != "" {
		filter = append(filter, bson.E{"name", name})
	}

	if realName := r.FormValue("real_name"); realName != "" {
		filter = append(filter, bson.E{"real_name", realName})
	}

	if email := r.FormValue("email"); email != "" {
		filter = append(filter, bson.E{"email", email})
	}

	cnt, _ := userModel.CountDocuments(context.Background(), filter)
	pg := pagination.Parse(r, cnt)

	opt := options.Find().SetSort(bson.D{
		{"create_time", -1},
	}).SetProjection(bson.D{
		{"password", 0},
		{"setting", 0},
	}).SetSkip(pg.SkipNum).SetLimit(pg.PageSize)

	items := []model.User{}
	if cur, err := userModel.Find(context.Background(), filter, opt); err == nil {
		cur.All(context.TODO(), &items)
	}

	rpr := response.ResponsePaginationResult{
		Total: cnt,
		Items: items,
	}

	response.JSON(rw, 0, rpr, "")
}

func (self *Service) Route_POST_Create(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var requestMUC predefined.RequestManageUserCreate
	if err := validator.FormStruct(&requestMUC, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	userModel := user.NewModel(self.M)
	if _, err := userModel.Create(context.Background(), &model.User{
		Name:        requestMUC.UserName,
		RealName:    requestMUC.RealName,
		Email:       requestMUC.Email,
		Password:    string(help.NewString(requestMUC.Password).Md5()),
		Description: requestMUC.Description,
		Enable:      requestMUC.Enable,
	}); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_POST_Update(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var requestMUU predefined.RequestManageUserUpdate
	if err := validator.FormStruct(&requestMUU, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	change := bson.D{
		{"name", requestMUU.UserName},
		{"real_name", requestMUU.RealName},
		{"email", requestMUU.Email},
		{"description", requestMUU.Description},
		{"enable", requestMUU.Enable},
	}

	userModel := user.NewModel(self.M)
	if ur, err := userModel.Set(r.Context(), userModel.FilterByID(requestMUU.UserID), change); err == nil && ur.ModifiedCount > 0 {
		response.JSON(rw, 0, nil, "")
	} else {
		log.Debug("error:", err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	}
}

func (self *Service) Route_POST_ResetPassword(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var requestMUID predefined.RequestManageUserID
	if err := validator.FormStruct(&requestMUID, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	password := help.NewString("").RandMixed(6)

	change := bson.D{
		{"password", help.NewString(password).Md5()},
	}

	userModel := user.NewModel(self.M)
	if ur, err := userModel.Set(r.Context(), userModel.FilterByID(requestMUID.UserID), change); err == nil && ur.ModifiedCount > 0 {
		response.JSON(rw, 0, password, "")
	} else {
		log.Debug("error:", err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	}
}

func (self *Service) Route_POST_Enable(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var requestMUE predefined.RequestManageUserEnable
	if err := validator.FormStruct(&requestMUE, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	userModel := user.NewModel(self.M)
	filter := userModel.FilterByID(requestMUE.UserID)

	change := bson.D{
		{"enable", requestMUE.Enable},
	}

	if _, err := userModel.Set(r.Context(), filter, change); err == nil {
		response.JSON(rw, 0, nil, "")
	} else {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	}
}

func (self *Service) Route_POST_Delete(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var requestMUD predefined.RequestManageUserDelete
	if err := validator.FormStruct(&requestMUD, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	userModel := user.NewModel(self.M)
	filter := userModel.FilterByID(requestMUD.UserID)

	change := bson.D{
		{"delete", requestMUD.Delete},
	}

	if _, err := userModel.Set(r.Context(), filter, change); err == nil {
		response.JSON(rw, 0, nil, "")
	} else {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	}
}

func (self *Service) Route_DELETE_Delete(rw http.ResponseWriter, r *http.Request) {
	vs, _ := request.ParseDeleteForm(r)

	var requestMUID predefined.RequestManageUserID
	if err := validator.FormStruct(&requestMUID, vs); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	userModel := user.NewModel(self.M)
	filter := userModel.FilterByID(requestMUID.UserID)
	if _, err := userModel.DeleteOne(r.Context(), filter); err == nil {
		response.JSON(rw, 0, nil, "")
	} else {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	}
}
