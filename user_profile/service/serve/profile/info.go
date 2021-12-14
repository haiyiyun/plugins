package profile

import (
	"net/http"

	"github.com/haiyiyun/plugins/user_profile/database/model"
	"github.com/haiyiyun/plugins/user_profile/database/model/profile"
	"github.com/haiyiyun/plugins/user_profile/predefined"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/utils/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (self *Service) Route_POST_Nickname(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var requestNn predefined.RequestServeNickname
	if err := validator.FormStruct(&requestNn, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.Set(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.nickname", requestNn.Nickname},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}

func (self *Service) Route_POST_Avatar(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var requestA predefined.RequestServeAvatar
	if err := validator.FormStruct(&requestA, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.Set(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.avatar", requestA.Avatar},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}

func (self *Service) Route_POST_Photos(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var requestP predefined.RequestServePhotos
	if err := validator.FormStruct(&requestP, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.Set(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.photos", requestP.Photos},
	}); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}

func (self *Service) Route_POST_Tags(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var requestT predefined.RequestServeTags
	if err := validator.FormStruct(&requestT, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.Set(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.tags", requestT.Tags},
	}); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}

func (self *Service) Route_POST_Education(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var pie model.ProfileInfoEducation
	if err := validator.FormStruct(&pie, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.Set(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.education", pie},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}

func (self *Service) Route_POST_Profession(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var pip model.ProfileInfoProfession
	if err := validator.FormStruct(&pip, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.Set(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.profession", pip},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}

func (self *Service) Route_POST_Contact(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var pic model.ProfileInfoContact
	if err := validator.FormStruct(&pic, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.Set(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.contact", pic},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}

func (self *Service) Route_POST_CoverImage(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var requestCI predefined.RequestServeCoverImage
	if err := validator.FormStruct(&requestCI, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.Set(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.cover.image", requestCI.Image},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}

func (self *Service) Route_POST_CoverVideo(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var requestCV predefined.RequestServeCoverVideo
	if err := validator.FormStruct(&requestCV, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.Set(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.cover.video", requestCV.Video},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}

func (self *Service) Route_POST_CoverVoice(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var requestCV predefined.RequestServeCoverVoice
	if err := validator.FormStruct(&requestCV, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.Set(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.cover.voice", requestCV.Voice},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}
