package user

import (
	"context"
	"net/http"
	"strconv"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/plugins/urbac/database/model"
	"github.com/haiyiyun/plugins/urbac/database/model/user"
	"github.com/haiyiyun/validator"

	"github.com/haiyiyun/utils/help"
	"github.com/haiyiyun/utils/http/pagination"
	"github.com/haiyiyun/utils/http/request"
	"github.com/haiyiyun/utils/http/response"
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

func (self *Service) Route_PUT_Create(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.FormValue("username")
	realName := r.FormValue("real_name")
	email := r.FormValue("email")
	enableStr := r.FormValue("enable")
	password := r.FormValue("password")
	description := r.FormValue("description")

	valid := &validator.Validation{}
	valid.Required(username).Key("username").Message("username字段为空")
	valid.Required(realName).Key("real_name").Message("real_name字段为空")
	valid.Required(email).Key("email").Message("email字段为空")
	valid.Required(enableStr).Key("enable").Message("enable字段为空")
	valid.Required(password).Key("password").Message("password字段为空")
	valid.Required(description).Key("description").Message("description字段为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	enable, _ := strconv.ParseBool(enableStr)

	userModel := user.NewModel(self.M)
	if _, err := userModel.Create(context.Background(), &model.User{
		Name:        username,
		RealName:    realName,
		Email:       email,
		Password:    string(help.NewString(password).Md5()),
		Description: description,
		Enable:      enable,
	}); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_POST_Update(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userIDHex := r.FormValue("user_id")
	username := r.FormValue("username")
	realName := r.FormValue("real_name")
	email := r.FormValue("email")
	enableStr := r.FormValue("enable")
	description := r.FormValue("description")

	valid := &validator.Validation{}
	valid.Required(userIDHex).Key("user_id").Message("user_id字段为空")
	valid.Required(username).Key("username").Message("username字段为空")
	valid.Required(realName).Key("real_name").Message("real_name字段为空")
	valid.Required(email).Key("email").Message("email字段为空")
	valid.Required(enableStr).Key("enable").Message("enable字段为空")
	valid.Required(description).Key("description").Message("description字段为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	userID, userIDErr := primitive.ObjectIDFromHex(userIDHex)
	if userIDErr != nil {
		response.JSON(rw, http.StatusBadRequest, nil, "")
		return
	}

	enable, _ := strconv.ParseBool(enableStr)

	change := bson.D{
		{"name", username},
		{"real_name", realName},
		{"email", email},
		{"description", description},
		{"enable", enable},
	}

	userModel := user.NewModel(self.M)
	if ur, err := userModel.Set(r.Context(), userModel.FilterByID(userID), change); err == nil && ur.ModifiedCount > 0 {
		response.JSON(rw, 0, nil, "")
	} else {
		log.Debug("error:", err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	}
}

func (self *Service) Route_POST_ResetPassword(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userIDHex := r.FormValue("user_id")

	valid := &validator.Validation{}

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	userID, userIDErr := primitive.ObjectIDFromHex(userIDHex)
	if userIDErr != nil {
		response.JSON(rw, http.StatusBadRequest, nil, "")
		return
	}

	password := help.NewString("").RandMixed(6)

	change := bson.D{
		{"password", help.NewString(password).Md5()},
	}

	userModel := user.NewModel(self.M)
	if ur, err := userModel.Set(r.Context(), userModel.FilterByID(userID), change); err == nil && ur.ModifiedCount > 0 {
		response.JSON(rw, 0, password, "")
	} else {
		log.Debug("error:", err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	}
}

func (self *Service) Route_POST_Enable(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userIDHex := r.FormValue("user_id")
	enableStr := r.FormValue("enable")
	valid := &validator.Validation{}
	valid.Required(userIDHex).Key("user_id").Message("user_id字段为空")
	valid.Required(enableStr).Key("enable").Message("enable字段为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	if userID, err := primitive.ObjectIDFromHex(userIDHex); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		userModel := user.NewModel(self.M)
		filter := userModel.FilterByID(userID)
		enable, _ := strconv.ParseBool(enableStr)

		change := bson.D{
			{"enable", enable},
		}

		if _, err = userModel.Set(r.Context(), filter, change); err == nil {
			response.JSON(rw, 0, nil, "")
		} else {
			response.JSON(rw, http.StatusBadRequest, nil, "")
		}
	}
}

func (self *Service) Route_POST_Delete(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userIDHex := r.FormValue("user_id")
	deleteStr := r.FormValue("delete")
	valid := &validator.Validation{}
	valid.Required(userIDHex).Key("user_id").Message("user_id字段为空")
	valid.Required(deleteStr).Key("delete").Message("delete字段为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	if userID, err := primitive.ObjectIDFromHex(userIDHex); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		userModel := user.NewModel(self.M)
		filter := userModel.FilterByID(userID)
		del, _ := strconv.ParseBool(deleteStr)

		change := bson.D{
			{"delete", del},
		}

		if _, err = userModel.Set(r.Context(), filter, change); err == nil {
			response.JSON(rw, 0, nil, "")
		} else {
			response.JSON(rw, http.StatusBadRequest, nil, "")
		}
	}
}

func (self *Service) Route_DELETE_Delete(rw http.ResponseWriter, r *http.Request) {
	vs, _ := request.ParseDeleteForm(r)

	userIDHex := vs.Get("user_id")

	valid := &validator.Validation{}
	valid.Required(userIDHex).Key("user_id").Message("user_id字段为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	if userID, err := primitive.ObjectIDFromHex(userIDHex); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		userModel := user.NewModel(self.M)
		filter := userModel.FilterByID(userID)
		if _, err = userModel.DeleteOne(r.Context(), filter); err == nil {
			response.JSON(rw, 0, nil, "")
		} else {
			response.JSON(rw, http.StatusBadRequest, nil, "")
		}
	}
}
