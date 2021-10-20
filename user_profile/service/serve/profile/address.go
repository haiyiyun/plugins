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

func (self *Service) Route_POST_CreateAddress(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	typStr := r.FormValue("type")
	typ, _ := strconv.Atoi(typStr)
	nation := r.FormValue("nation")
	province := r.FormValue("province")
	city := r.FormValue("city")
	district := r.FormValue("district")
	address := r.FormValue("address")

	valid := validator.Validation{}
	valid.Digital(typStr).Key("type").Message("type必须为数字")
	valid.Have(typ, predefined.ProfileInfoAddressTypeHometown, predefined.ProfileInfoAddressTypeResidence).Key("type").Message("请提供支持的type")
	valid.Required(nation).Key("nation").Message("nation不能为空")
	valid.Required(province).Key("province").Message("province不能为空")
	valid.Required(city).Key("city").Message("city不能为空")

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
		{"info.address", model.ProfileInfoAddress{
			Type:     typ,
			Nation:   nation,
			Province: province,
			City:     city,
			District: district,
			Address:  address,
		}},
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
	typStr := r.FormValue("type")
	typ, _ := strconv.Atoi(typStr)
	nation := r.FormValue("nation")
	province := r.FormValue("province")
	city := r.FormValue("city")
	district := r.FormValue("district")
	address := r.FormValue("address")

	valid := validator.Validation{}
	valid.Digital(typStr).Key("type").Message("type必须为数字")
	valid.Have(typ, predefined.ProfileInfoAddressTypeHometown, predefined.ProfileInfoAddressTypeResidence).Key("type").Message("请提供支持的type")
	valid.Required(nation).Key("nation").Message("nation不能为空")
	valid.Required(province).Key("province").Message("province不能为空")
	valid.Required(city).Key("city").Message("city不能为空")

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
		{"info.address.type", typ},
	}...)

	if ur, err := profileModel.Set(r.Context(), filter, bson.D{
		{"info.address.$.nation", nation},
		{"info.address.$.province", province},
		{"info.address.$.city", city},
		{"info.address.$.district", district},
		{"info.address.$.address", address},
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
	typStr := r.FormValue("type")
	typ, _ := strconv.Atoi(typStr)
	nation := r.FormValue("nation")
	province := r.FormValue("province")
	city := r.FormValue("city")
	district := r.FormValue("district")
	address := r.FormValue("address")

	valid := validator.Validation{}
	valid.Digital(typStr).Key("type").Message("type必须为数字")
	valid.Have(typ, predefined.ProfileInfoAddressTypeHometown, predefined.ProfileInfoAddressTypeResidence).Key("type").Message("请提供支持的type")
	valid.Required(nation).Key("nation").Message("nation不能为空")
	valid.Required(province).Key("province").Message("province不能为空")
	valid.Required(city).Key("city").Message("city不能为空")

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
		{"info.address", model.ProfileInfoAddress{
			Type:     typ,
			Nation:   nation,
			Province: province,
			City:     city,
			District: district,
			Address:  address,
		}},
	}); err != nil || ur.ModifiedCount == 0 {
		log.Error(err)
		response.JSON(rw, http.StatusBadRequest, nil, "")
	} else {
		response.JSON(rw, 0, nil, "")
		return
	}
}
