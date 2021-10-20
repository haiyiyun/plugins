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

func (self *Service) Route_POST_CreateQuestion(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	typStr := r.FormValue("type")
	typ, _ := strconv.Atoi(typStr)
	question := r.FormValue("question")

	valid := validator.Validation{}
	valid.Digital(typStr).Key("type").Message("type必须为数字")
	valid.Have(typ, predefined.ProfileQuestionTypeValues).Key("type").Message("请提供支持的type")

	valid.Required(question).Key("question").Message("question不能为空")

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
	if ur, err := profileModel.AddToSet(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"questions", model.ProfileQuestion{
			Type:     typ,
			Question: question,
		}},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}

func (self *Service) Route_POST_UpdateQuestion(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	typStr := r.FormValue("type")
	typ, _ := strconv.Atoi(typStr)
	question := r.FormValue("question")

	valid := validator.Validation{}
	valid.Digital(typStr).Key("type").Message("type必须为数字")
	valid.Have(typ, predefined.ProfileQuestionTypeValues).Key("type").Message("请提供支持的type")

	valid.Required(question).Key("question").Message("question不能为空")

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
		{"questions.type", typ},
	}...)
	// opt := options.Update().SetArrayFilters(options.ArrayFilters{
	// 	Filters: []interface{}{
	// 		bson.D{
	// 			{"elem.type", typ},
	// 		},
	// 	},
	// })

	if ur, err := profileModel.Set(r.Context(), filter, bson.D{
		// {"questions.$[elem].question", question},
		{"questions.$.question", question},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}

func (self *Service) Route_POST_DeleteQuestion(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	typStr := r.FormValue("type")
	typ, _ := strconv.Atoi(typStr)
	question := r.FormValue("question")

	valid := validator.Validation{}
	valid.Digital(typStr).Key("type").Message("type必须为数字")
	valid.Have(typ, predefined.ProfileQuestionTypeValues).Key("type").Message("请提供支持的type")

	valid.Required(question).Key("question").Message("question不能为空")

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
	if ur, err := profileModel.Pull(r.Context(), profileModel.FilterByID(userID), bson.D{
		{"questions", model.ProfileQuestion{
			Type:     typ,
			Question: question,
		}},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}
