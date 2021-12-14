package role

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/plugins/urbac/database/model"
	"github.com/haiyiyun/plugins/urbac/database/model/role"
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

	userHexs := r.Form["users[]"]

	var requestMR predefined.RequestManageRole
	if err := validator.FormStruct(&requestMR, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	roleApps := map[string]model.Application{}
	jsonBlob := json.RawMessage(requestMR.Applications)
	if err := json.Unmarshal(jsonBlob, &roleApps); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, "")
		return
	}

	roleModel := role.NewModel(self.M)
	if _, err := roleModel.Create(context.Background(), &model.Role{
		Name:   requestMR.Name,
		Enable: requestMR.Enable,
		Users:  help.NewSlice(userHexs).ConvObjectID(),
		Right: model.RoleRight{
			Scope:        requestMR.Scope,
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

	userHexs := r.Form["users[]"]

	var requestMR predefined.RequestManageRole
	if err := validator.FormStruct(&requestMR, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	roleApps := map[string]model.Application{}
	jsonBlob := json.RawMessage(requestMR.Applications)
	if err := json.Unmarshal(jsonBlob, &roleApps); err != nil {
		log.Debug("error:", err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
		return
	}

	change := bson.D{
		{"users", help.NewSlice(userHexs).ConvObjectID()},
		{"right.scope", requestMR.Scope},
		{"right.applications", roleApps},
		{"enable", requestMR.Enable},
	}

	roleModel := role.NewModel(self.M)
	if ur, err := roleModel.Set(r.Context(), roleModel.FilterByName(requestMR.Name), change); err == nil && ur.ModifiedCount > 0 {
		response.JSON(rw, 0, nil, "")
	} else {
		log.Debug("error:", err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	}
}

func (self *Service) Route_POST_Enable(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var requestMRE predefined.RequestManageRoleEnable
	if err := validator.FormStruct(&requestMRE, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	roleModel := role.NewModel(self.M)
	filter := roleModel.FilterByID(requestMRE.ObjectID)

	change := bson.D{
		{"enable", requestMRE.Enable},
	}

	if _, err := roleModel.Set(r.Context(), filter, change); err == nil {
		response.JSON(rw, 0, nil, "")
	} else {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	}
}

func (self *Service) Route_POST_Delete(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var requestMRD predefined.RequestManageRoleDelete
	if err := validator.FormStruct(&requestMRD, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	roleModel := role.NewModel(self.M)
	filter := roleModel.FilterByID(requestMRD.ObjectID)

	change := bson.D{
		{"delete", requestMRD.Delete},
	}

	if _, err := roleModel.Set(r.Context(), filter, change); err == nil {
		response.JSON(rw, 0, nil, "")
	} else {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	}
}

func (self *Service) Route_DELETE_Delete(rw http.ResponseWriter, r *http.Request) {
	vs, _ := request.ParseDeleteForm(r)

	var requestMO predefined.RequestManageObjectID
	if err := validator.FormStruct(&requestMO, vs); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	roleModel := role.NewModel(self.M)
	filter := roleModel.FilterByID(requestMO.ObjectID)
	if _, err := roleModel.DeleteOne(r.Context(), filter); err == nil {
		response.JSON(rw, 0, nil, "")
	} else {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	}
}
