package profile

import (
	"net/http"

	"github.com/haiyiyun/plugins/user_profile/database/model"
	"github.com/haiyiyun/plugins/user_profile/database/model/profile"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/utils/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (self *Service) Route_POST_CreateIntroduction(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var pii model.ProfileInfoIntroduction
	if err := validator.FormStruct(&pii, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.AddToSet(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.introduction", pii},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}

func (self *Service) Route_POST_UpdateIntroduction(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var pii model.ProfileInfoIntroduction
	if err := validator.FormStruct(&pii, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	profileModel := profile.NewModel(self.M)
	filter := profileModel.FilterByID(userID)
	filter = append(filter, bson.D{
		{"info.introduction.type", pii.Type},
	}...)

	if ur, err := profileModel.Set(r.Context(), filter, bson.D{
		{"info.introduction.$.introduction", pii.Introduction},
		{"info.introduction.$.photos", pii.Photos},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}

func (self *Service) Route_POST_DeleteIntroduction(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var pii model.ProfileInfoIntroduction
	if err := validator.FormStruct(&pii, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.Pull(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.introduction", pii},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}
