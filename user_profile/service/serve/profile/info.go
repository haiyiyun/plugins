package profile

import (
	"net/http"
	"strconv"

	"github.com/haiyiyun/plugins/user_profile/database/model"
	"github.com/haiyiyun/plugins/user_profile/database/model/profile"
	"github.com/haiyiyun/plugins/user_profile/predefined"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (self *Service) Route_POST_Nickname(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nickname := r.FormValue("nickname")

	valid := validator.Validation{}
	valid.Required(nickname).Key("nickname").Message("nickname不能为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.Set(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.nickname", nickname},
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
	avatar := r.FormValue("avatar")

	valid := validator.Validation{}
	valid.Required(avatar).Key("avatar").Message("avatar不能为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.Set(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.avatar", avatar},
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
	photos := r.Form["photo[]"]

	valid := validator.Validation{}
	valid.Required(photos).Key("photos").Message("photos不能为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.Set(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.photos", photos},
	}); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}

func (self *Service) Route_POST_Tags(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tags := r.Form["tag[]"]

	valid := validator.Validation{}
	valid.Required(tags).Key("tag").Message("tag不能为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.Set(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.tags", tags},
	}); err != nil || ur.ModifiedCount == 0 {
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}

func (self *Service) Route_POST_Education(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	highestDegree := r.FormValue("highest_degree")
	graduatedCollege := r.FormValue("graduated_college")

	valid := validator.Validation{}
	valid.Required(highestDegree).Key("highest_degree").Message("highest_degree不能为空")
	valid.Required(graduatedCollege).Key("graduated_college").Message("graduated_college不能为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.Set(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.education", model.ProfileInfoEducation{
			HighestDegree:    highestDegree,
			GraduatedCollege: graduatedCollege,
		}},
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
	company := r.FormValue("company")
	position := r.FormValue("position")
	annualIncomeStr := r.FormValue("annual_income")
	annualIncome, _ := strconv.Atoi(annualIncomeStr)

	valid := validator.Validation{}
	valid.Required(company).Key("company").Message("company不能为空")
	valid.Required(position).Key("position").Message("position不能为空")
	valid.Digital(annualIncomeStr).Key("annual_income").Message("annual_income必须数字")
	valid.Have(annualIncome,
		predefined.ProfileInfoProfessionAnnualIncome5_15,
		predefined.ProfileInfoProfessionAnnualIncome15_30,
		predefined.ProfileInfoProfessionAnnualIncome30_50,
		predefined.ProfileInfoProfessionAnnualIncome50_100,
		predefined.ProfileInfoProfessionAnnualIncome100_500,
		predefined.ProfileInfoProfessionAnnualIncome500_,
	).Key("annual_income").Message("请提供支持的annual_income")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.Set(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.profession", model.ProfileInfoProfession{
			Company:      company,
			Position:     position,
			AnnualIncome: annualIncome,
		}},
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
	phoneNumber := r.FormValue("phone_number")
	email := r.FormValue("email")

	valid := validator.Validation{}
	valid.ChinaMobile(phoneNumber).Key("phone_number").Message("phone_number必须为正确的手机号")
	valid.Email(email).Key("email").Message("email必须为正确的邮箱地址")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.Set(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.contact", model.ProfileInfoContact{
			PhoneNumber: phoneNumber,
			Email:       email,
		}},
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
	image := r.FormValue("image")

	valid := validator.Validation{}
	valid.Required(image).Key("image").Message("image不能为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.Set(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.cover.image", image},
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
	video := r.FormValue("video")

	valid := validator.Validation{}
	valid.Required(video).Key("video").Message("video不能为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.Set(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.cover.video", video},
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
	voice := r.FormValue("voice")

	valid := validator.Validation{}
	valid.Required(voice).Key("voice").Message("voice不能为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	profileModel := profile.NewModel(self.M)
	if ur, err := profileModel.Set(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"info.cover.voice", voice},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}
