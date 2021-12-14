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

func (self *Service) Route_POST_CreateAddress(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var pia model.ProfileInfoAddress
	if err := validator.FormStruct(&pia, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.AddToSet(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.address", pia},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}

func (self *Service) Route_POST_UpdateAddress(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var pia model.ProfileInfoAddress
	if err := validator.FormStruct(&pia, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	profileModel := profile.NewModel(self.M)
	filter := profileModel.FilterByID(userID)
	filter = append(filter, bson.D{
		{"info.address.type", pia.Type},
	}...)

	if ur, err := profileModel.Set(r.Context(), filter, bson.D{
		{"info.address.$.nation", pia.Nation},
		{"info.address.$.province", pia.Province},
		{"info.address.$.city", pia.City},
		{"info.address.$.district", pia.District},
		{"info.address.$.address", pia.Address},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}

func (self *Service) Route_POST_DeleteAddress(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var pia model.ProfileInfoAddress
	if err := validator.FormStruct(&pia, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.Pull(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.address", pia},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}
