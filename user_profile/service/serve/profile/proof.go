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

func (self *Service) Route_POST_ProofIdentityCard(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("id")
	realName := r.FormValue("real_name")
	images := r.Form["images[]"]

	valid := validator.Validation{}
	valid.Required(id).Key("id").Message("id不能为空")
	valid.Required(realName).Key("real_name").Message("real_name不能为空")
	valid.Required(images).Key("images").Message("images不能为空")

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

	filter := profileModel.FilterByID(userID)
	filter = append(filter, bson.D{
		{
			"$or", []bson.D{
				{{"proof.identity_card.verified", bson.D{
					{"$exists", false},
				}}},
				{{"proof.identity_card.verified", false}},
			},
		},
	}...)

	if ur, err := profileModel.Set(r.Context(), filter, bson.D{
		{"proof.identity_card", model.ProfileProofIdentityCard{
			ID:       id,
			RealName: realName,
			Images:   images,
		}},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}

func (self *Service) Route_POST_ProofProfession(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	typStr := r.FormValue("type")
	typ, _ := strconv.Atoi(typStr)
	companyName := r.FormValue("company_name")
	showName := r.FormValue("show_name")
	schoolName := r.FormValue("school_name")
	degree := r.FormValue("degree")
	images := r.Form["images[]"]

	valid := validator.Validation{}
	valid.Digital(typStr).Key("type").Message("type必须为数字")
	valid.Have(typ,
		predefined.ProfileProofProfessionTypeCompanySocialSecurity,
		predefined.ProfileProofProfessionTypeCompanyEnterpriseOfficeSoftware,
		predefined.ProfileProofProfessionTypeCompanyLicense,
		predefined.ProfileProofProfessionTypeCompanyWorkPermit,
		predefined.ProfileProofProfessionTypeCompanyPaySlip,
		predefined.ProfileProofProfessionTypeCompanyOffer,
		predefined.ProfileProofProfessionTypeStudent,
	).Key("type").Message("必须支持的type")

	if typ == predefined.ProfileProofProfessionTypeStudent {
		valid.Required(schoolName).Key("school_name").Message("school_name不能为空")
		valid.Required(degree).Key("degree").Message("degree不能为空")
	} else {
		valid.Required(companyName).Key("company_name").Message("company_name不能为空")
		valid.Required(showName).Key("show_name").Message("show_name不能为空")
	}

	valid.Required(images).Key("images").Message("images不能为空")

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

	filter := profileModel.FilterByID(userID)
	filter = append(filter, bson.D{
		{
			"$or", []bson.D{
				{{"proof.profession.verified", bson.D{
					{"$exists", false},
				}}},
				{{"proof.profession.verified", false}},
			},
		},
	}...)

	setProfileProofProfession := model.ProfileProofProfession{
		Type: typ,
	}

	if typ == predefined.ProfileProofProfessionTypeStudent {
		setProfileProofProfession.Student = model.ProfileProofStudent{
			SchoolName: schoolName,
			Degree:     degree,
			Images:     images,
		}
	} else {
		setProfileProofProfession.Company = model.ProfileProofCompany{
			CompanyName: companyName,
			ShowName:    showName,
			Images:      images,
		}
	}

	if ur, err := profileModel.Set(r.Context(), filter, bson.D{
		{"proof.profession", setProfileProofProfession},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}

func (self *Service) Route_POST_ProofEducation(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	typStr := r.FormValue("type")
	typ, _ := strconv.Atoi(typStr)
	id := r.FormValue("id")
	collegeName := r.FormValue("college_name")
	degree := r.FormValue("degree")
	images := r.Form["images[]"]
	place := r.FormValue("place")
	year := r.FormValue("year")

	valid := validator.Validation{}
	valid.Digital(typStr).Key("type").Message("type必须为数字")
	valid.Have(typ,
		predefined.ProfileProofEducationTypeCHSI,
		predefined.ProfileProofEducationTypeDiplomaImage,
		predefined.ProfileProofEducationTypeDiplomaID,
		predefined.ProfileProofEducationTypeCSCSE,
		predefined.ProfileProofEducationTypeOldCSCSE,
	).Key("type").Message("必须支持的type")

	valid.Required(id).Key("id").Message("id不能为空")
	valid.Required(collegeName).Key("college_name").Message("college_name不能为空")
	valid.Required(degree).Key("degree").Message("degree不能为空")

	if typ == predefined.ProfileProofEducationTypeDiplomaImage {
		valid.Required(images).Key("images").Message("images不能为空")
	}

	if typ == predefined.ProfileProofEducationTypeOldCSCSE {
		valid.Required(place).Key("place").Message("place不能为空")
		valid.Required(year).Key("year").Message("year不能为空")
	}

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

	filter := profileModel.FilterByID(userID)
	filter = append(filter, bson.D{
		{
			"$or", []bson.D{
				{{"proof.education.verified", bson.D{
					{"$exists", false},
				}}},
				{{"proof.education.verified", false}},
			},
		},
	}...)

	setProfileProofEducation := model.ProfileProofEducation{
		Type:        typ,
		ID:          id,
		CollegeName: collegeName,
		Degree:      degree,
	}

	if typ == predefined.ProfileProofEducationTypeDiplomaImage {
		setProfileProofEducation.Images = images
	}

	if typ == predefined.ProfileProofEducationTypeOldCSCSE {
		setProfileProofEducation.Place = place
		setProfileProofEducation.Year = year
	}

	if ur, err := profileModel.Set(r.Context(), filter, bson.D{
		{"proof.education", setProfileProofEducation},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}
