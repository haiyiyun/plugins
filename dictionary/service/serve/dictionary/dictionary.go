package dictionary

import (
	"net/http"
	"strconv"

	"github.com/haiyiyun/plugins/dictionary/database/model/dictionary"
	"github.com/haiyiyun/plugins/dictionary/predefined"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/validator"
)

func (self *Service) Route_GET_Lookup(rw http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	structrue := r.URL.Query().Get("structure")
	want := r.URL.Query().Get("want")
	flushStr := r.URL.Query().Get("flush")

	valid := &validator.Validation{}
	valid.Required(key).Key("key").Message("key字段为空")
	valid.Have(structrue, predefined.DictionaryLookupStructureSlice, predefined.DictionaryLookupStructureMap).Key("structrue").Message("structrue字段错误")
	valid.Have(want, predefined.DictionaryLookupWantLable, predefined.DictionaryLookupWantValue, predefined.DictionaryLookupWantLableAndValue).Key("want").Message("want字段错误")
	valid.Required(flushStr).Key("flush").Message("flush字段为空")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	flush, _ := strconv.ParseBool(flushStr)

	cacheKey := "dictionary." + key + "." + structrue + "." + want
	if dict, found := self.Cache.Get(cacheKey); !flush && found {
		response.JSON(rw, 0, dict, "")
	} else {
		dictModel := dictionary.NewModel(self.M)

		if dict := dictModel.Lookup(key, structrue, want); dict != nil {
			self.Cache.Set(cacheKey, dict, 0)
			response.JSON(rw, 0, dict, "")
			return
		}
	}

	response.JSON(rw, http.StatusBadRequest, nil, "")
}
