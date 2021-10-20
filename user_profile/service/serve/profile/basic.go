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

func (self *Service) Route_POST_Basic(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	birghYearStr := r.FormValue("birth_year")
	birghYear, _ := strconv.Atoi(birghYearStr)
	birthMonthStr := r.FormValue("birth_month")
	birthMonth, _ := strconv.Atoi(birthMonthStr)
	birthDayStr := r.FormValue("birth_day")
	birthDay, _ := strconv.Atoi(birthDayStr)
	sexStr := r.FormValue("sex")
	sex, _ := strconv.Atoi(sexStr)
	heightStr := r.FormValue("height")
	height, _ := strconv.Atoi(heightStr)
	weightStr := r.FormValue("weight")
	weight, _ := strconv.Atoi(weightStr)
	marriageStr := r.FormValue("marriage")
	marriage, _ := strconv.Atoi(marriageStr)
	constellationStr := r.FormValue("constellation")
	constellation, _ := strconv.Atoi(constellationStr)

	valid := validator.Validation{}
	valid.Digital(sexStr).Key("sexStr").Message("sex必须为数字")
	valid.Have(sex, predefined.ProfileInfoSexFemale, predefined.ProfileInfoSexMale).Key("sex").Message("请提供支持的sex")
	valid.Digital(birghYearStr).Key("birth_year").Message("birth_year必须为数字")
	valid.Digital(birthMonthStr).Key("birth_month_str").Message("birth_month必须为数字")
	valid.Range(birthMonth, 0, 11).Key("birth_month").Message("birth_month必须在0-11之间")
	valid.Digital(birthDayStr).Key("birth_day_str").Message("birth_day必须为数字")
	valid.Range(birthDay, 0, 30).Key("birth_day").Message("birth_day必须在0-30之间")
	valid.Digital(heightStr).Key("height_str").Message("height必须为数字")
	valid.Range(height, 140, 250).Key("height").Message("height必须在140-250之间")
	valid.Digital(weightStr).Key("weight_str").Message("weight必须为数字")
	valid.Range(weight, 50, 250).Key("weight").Message("weight必须在50-250之间")
	valid.Digital(marriageStr).Key("marriage_str").Message("marriage必须为数字")
	valid.Have(marriage,
		predefined.ProfileInfoBasicMarriageUnmarried,
		predefined.ProfileInfoBasicMarriageMarried,
		predefined.ProfileInfoBasicMarriageDivorced,
	).Key("marriage").Message("请提供支持的marriage")
	valid.Digital(constellationStr).Key("constellatione_str").Message("constellation必须为数字")
	valid.Have(constellation,
		predefined.ProfileInfoConstellationAries,
		predefined.ProfileInfoConstellationTaurus,
		predefined.ProfileInfoConstellationGemini,
		predefined.ProfileInfoConstellationCancer,
		predefined.ProfileInfoConstellationLeo,
		predefined.ProfileInfoConstellationVirgo,
		predefined.ProfileInfoConstellationLibra,
		predefined.ProfileInfoConstellationScorpio,
		predefined.ProfileInfoConstellationSagittarius,
		predefined.ProfileInfoConstellationCapricorn,
		predefined.ProfileInfoConstellationAquarius,
		predefined.ProfileInfoConstellationPisces,
	).Key("constellatione").Message("请提供支持的constellation")

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
		{"info.basic", model.ProfileInfoBasic{
			Sex: sex,
			Birth: model.ProfileInfoBirth{
				Year:  birghYear,
				Month: birthMonth,
				Day:   birthDay,
			},
			Height:        height,
			Weight:        weight,
			Marriage:      marriage,
			Constellation: constellation,
		}},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}
