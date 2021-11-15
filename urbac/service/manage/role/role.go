package role

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/plugins/urbac/database/model"
	"github.com/haiyiyun/plugins/urbac/database/model/role"
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
	roleModel := role.NewModel(self.M)
	filter := bson.D{}

	rightScopes := r.Form["right.scope[]"]
	if len(rightScopes) > 0 {
		filter = append(filter, bson.E{
			"right.scope", bson.D{
				{"$in", help.NewSlice(rightScopes).ConvInt()},
			},
		})
	}

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

	if userIDHex := r.FormValue("user_id"); userIDHex != "" {
		if userID, err := primitive.ObjectIDFromHex(userIDHex); err == nil {
			filter = append(filter, bson.E{"users", userID})
		}
	}

	if name := r.FormValue("name"); name != "" {
		filter = append(filter, bson.E{"name", name})
	}

	cnt, _ := roleModel.CountDocuments(r.Context(), filter)
	pg := pagination.Parse(r, cnt)

	opt := options.Find().SetSort(bson.D{
		{"create_time", -1},
	}).SetProjection(bson.D{}).SetSkip(pg.SkipNum).SetLimit(pg.PageSize)

	items := []model.Role{}
	if cur, err := roleModel.Find(r.Context(), filter, opt); err == nil {
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

	name := r.FormValue("name")
	scopeStr := r.FormValue("scope")
	enableStr := r.FormValue("enable")
	applications := r.FormValue("applications")
	userHexs := r.Form["users[]"]

	valid := &validator.Validation{}
	valid.Required(name).Key("name").Message("name字段为空")
	valid.Required(scopeStr).Key("scope").Message("scope字段为空")
	valid.Required(enableStr).Key("enable").Message("enable字段为空")
	valid.Required(applications).Key("applications").Message("applications字段为空")
	// valid.Required(userHexs).Key("users").Message("users字段为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	roleApps := map[string]model.Application{}
	jsonBlob := json.RawMessage(applications)
	if err := json.Unmarshal(jsonBlob, &roleApps); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, "")
		return
	}

	scope, _ := strconv.Atoi(scopeStr)
	enable, _ := strconv.ParseBool(enableStr)

	roleModel := role.NewModel(self.M)
	if _, err := roleModel.Create(context.Background(), &model.Role{
		Name:   name,
		Enable: enable,
		Users:  help.NewSlice(userHexs).ConvObjectID(),
		Right: model.RoleRight{
			Scope:        scope,
			Applications: roleApps,
		},
	}); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}

func (self *Service) Route_POST_Update(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	name := r.FormValue("name")
	scopeStr := r.FormValue("scope")
	enableStr := r.FormValue("enable")
	applications := r.FormValue("applications")
	userHexs := r.Form["users[]"]

	valid := &validator.Validation{}
	valid.Required(name).Key("name").Message("name字段为空")
	valid.Required(scopeStr).Key("scope").Message("scope字段为空")
	valid.Required(enableStr).Key("enable").Message("enable字段为空")
	valid.Required(applications).Key("applications").Message("applications字段为空")
	// valid.Required(userHexs).Key("users").Message("users字段为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	roleApps := map[string]model.Application{}
	jsonBlob := json.RawMessage(applications)
	if err := json.Unmarshal(jsonBlob, &roleApps); err != nil {
		log.Debug("error:", err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
		return
	}

	scope, _ := strconv.Atoi(scopeStr)
	enable, _ := strconv.ParseBool(enableStr)

	change := bson.D{
		{"users", help.NewSlice(userHexs).ConvObjectID()},
		{"right.scope", scope},
		{"right.applications", roleApps},
		{"enable", enable},
	}

	roleModel := role.NewModel(self.M)
	if ur, err := roleModel.Set(r.Context(), roleModel.FilterByName(name), change); err == nil && ur.ModifiedCount > 0 {
		response.JSON(rw, 0, nil, "")
	} else {
		log.Debug("error:", err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	}
}

func (self *Service) Route_POST_Enable(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	roleIDHex := r.FormValue("_id")
	enableStr := r.FormValue("enable")
	valid := &validator.Validation{}
	valid.Required(roleIDHex).Key("_id").Message("_id字段为空")
	valid.Required(enableStr).Key("enable").Message("enable字段为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	if roleID, err := primitive.ObjectIDFromHex(roleIDHex); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		roleModel := role.NewModel(self.M)
		filter := roleModel.FilterByID(roleID)
		enable, _ := strconv.ParseBool(enableStr)

		change := bson.D{
			{"enable", enable},
		}

		if _, err = roleModel.Set(r.Context(), filter, change); err == nil {
			response.JSON(rw, 0, nil, "")
		} else {
			response.JSON(rw, http.StatusBadRequest, nil, "")
		}
	}
}

func (self *Service) Route_POST_Delete(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	roleIDHex := r.FormValue("_id")
	deleteStr := r.FormValue("delete")
	valid := &validator.Validation{}
	valid.Required(roleIDHex).Key("_id").Message("_id字段为空")
	valid.Required(deleteStr).Key("delete").Message("delete字段为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	if roleID, err := primitive.ObjectIDFromHex(roleIDHex); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		roleModel := role.NewModel(self.M)
		filter := roleModel.FilterByID(roleID)
		del, _ := strconv.ParseBool(deleteStr)

		change := bson.D{
			{"delete", del},
		}

		if _, err = roleModel.Set(r.Context(), filter, change); err == nil {
			response.JSON(rw, 0, nil, "")
		} else {
			response.JSON(rw, http.StatusBadRequest, nil, "")
		}
	}
}

func (self *Service) Route_DELETE_Delete(rw http.ResponseWriter, r *http.Request) {
	vs, _ := request.ParseDeleteForm(r)

	roleIDHex := vs.Get("_id")

	valid := &validator.Validation{}
	valid.Required(roleIDHex).Key("_id").Message("_id字段为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	if roleID, err := primitive.ObjectIDFromHex(roleIDHex); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		roleModel := role.NewModel(self.M)
		filter := roleModel.FilterByID(roleID)
		if _, err = roleModel.DeleteOne(r.Context(), filter); err == nil {
			response.JSON(rw, 0, nil, "")
		} else {
			response.JSON(rw, http.StatusBadRequest, nil, "")
		}
	}
}
