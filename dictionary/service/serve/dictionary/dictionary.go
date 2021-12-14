package dictionary

import (
	"net/http"

	"github.com/haiyiyun/plugins/dictionary/database/model/dictionary"
	"github.com/haiyiyun/plugins/dictionary/predefined"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/utils/validator"
)

func (self *Service) Route_GET_Lookup(rw http.ResponseWriter, r *http.Request) {
	var requestLookup predefined.RequestServeLookup
	if err := validator.FormStruct(&requestLookup, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	cacheKey := "dictionary." + requestLookup.Key + "." + requestLookup.Structure + "." + requestLookup.Want
	if dict, found := self.Cache.Get(cacheKey); !requestLookup.Flush && found {
		response.JSON(rw, 0, dict, "")
		return
	} else {
		dictModel := dictionary.NewModel(self.M)

		if dict := dictModel.Lookup(requestLookup.Key, requestLookup.Structure, requestLookup.Want); dict != nil {
			self.Cache.Set(cacheKey, dict, 0)
			response.JSON(rw, 0, dict, "")
			return
		}
	}

	response.JSON(rw, http.StatusBadRequest, nil, "")
}
