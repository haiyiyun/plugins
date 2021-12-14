package visitor

import (
	"net/http"

	"github.com/haiyiyun/plugins/user_relationship/database/model"
	"github.com/haiyiyun/plugins/user_relationship/database/model/visitor"
	"github.com/haiyiyun/plugins/user_relationship/predefined"

	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/utils/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (self *Service) Route_POST_Visitor(rw http.ResponseWriter, r *http.Request) {
	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	r.ParseForm()

	var requestV predefined.RequestServeVisitor
	if err := validator.FormStruct(&requestV, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	v := model.Visitor{
		OwnerUserID:   requestV.OwnerUserID,
		VisitorUserID: userID,
	}

	visitorModel := visitor.NewModel(self.M)
	if _, err := visitorModel.Create(r.Context(), v); err != nil {
		response.JSON(rw, http.StatusServiceUnavailable, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
	}
}
