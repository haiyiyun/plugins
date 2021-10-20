package profile

import (
	"net/http"

	"github.com/haiyiyun/plugins/user_profile/database/model/profile"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/haiyiyun/utils/http/response"
)

func (self *Service) Route_GET_Profile(rw http.ResponseWriter, r *http.Request) {
	if userID := self.GetUserID(r); userID != primitive.NilObjectID {
		profileModel := profile.NewModel(self.M)

		if pf, err := profileModel.GetInfo(userID); err == nil {
			pf["user_id"] = pf["_id"]
			delete(pf, "_id")
			response.JSON(rw, 0, pf, "")
			return
		}
	}

	response.JSON(rw, http.StatusNotFound, nil, "")
}
