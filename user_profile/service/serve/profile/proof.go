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

func (self *Service) Route_POST_ProofIdentityCard(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var ppidc model.ProfileProofIdentityCard
	if err := validator.FormStruct(&ppidc, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
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
		{"proof.identity_card", ppidc},
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

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var ppp model.ProfileProofProfession
	if err := validator.FormStruct(&ppp, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
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

	if ur, err := profileModel.Set(r.Context(), filter, bson.D{
		{"proof.profession", ppp},
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

	userID := self.GetUserID(r)
	if userID == primitive.NilObjectID {
		response.JSON(rw, http.StatusUnauthorized, nil, "")
		return
	}

	var ppe model.ProfileProofEducation
	if err := validator.FormStruct(&ppe, r.Form); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
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

	if ur, err := profileModel.Set(r.Context(), filter, bson.D{
		{"proof.education", ppe},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}
