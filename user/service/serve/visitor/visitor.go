package visitor

import (
	"net/http"

	"github.com/haiyiyun/plugins/user/database/model"
	"github.com/haiyiyun/plugins/user/database/model/visitor"

	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (self *Service) Route_POST_Visitor(rw http.ResponseWriter, r *http.Request) {
	u, found := self.GetUserInfo(r)
	if !found {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	ownerUserIDStr := r.FormValue("owner_user_id")
	ownerUserID, _ := primitive.ObjectIDFromHex(ownerUserIDStr)

	valid := validator.Validation{}
	valid.BsonObjectID(ownerUserIDStr).Key("owner_user_id").Message("owner_user_id必须支持的格式")
	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, "")
		return
	}

	v := model.Visitor{
		OwnerUserID:   ownerUserID,
		VisitorUserID: u.ID,
	}

	visitorModel := visitor.NewModel(self.M)
	if _, err := visitorModel.Create(r.Context(), v); err != nil {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}
