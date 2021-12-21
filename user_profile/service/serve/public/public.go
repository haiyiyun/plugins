package public

import (
	"net/http"

	"github.com/haiyiyun/plugins/user_profile/database/model/profile"
	"github.com/haiyiyun/plugins/user_profile/predefined"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/utils/validator"
)

func (self *Service) Route_GET_NicknameAvatar(rw http.ResponseWriter, r *http.Request) {
	var requestUID predefined.RequestServeUserID
	if err := validator.FormStruct(&requestUID, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	profileModel := profile.NewModel(self.M)
	if pf, err := profileModel.GetNickNameAndAvatar(requestUID.UserID); err == nil {
		response.JSON(rw, 0, pf, "")
		return
	}

	response.JSON(rw, http.StatusNotFound, nil, "")
}
