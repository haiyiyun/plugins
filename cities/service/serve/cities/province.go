package cities

import (
	"net/http"

	"github.com/haiyiyun/plugins/cities/database/model"
	"github.com/haiyiyun/plugins/cities/database/model/province"
	"github.com/haiyiyun/utils/http/response"
	"go.mongodb.org/mongo-driver/bson"
)

func (self *Service) Route_GET_Province(rw http.ResponseWriter, req *http.Request) {
	cacheKey := "cities.province"
	provinces := []model.Province{}

	if found, _ := self.Cache.Get(cacheKey, &provinces); found {
		response.JSON(rw, 0, provinces, "")
		return
	} else {
		provinceModel := province.NewModel(self.M)
		if cur, err := provinceModel.Find(req.Context(), bson.D{}); err == nil {
			if err := cur.All(req.Context(), &provinces); err == nil {
				self.Cache.Set(cacheKey, provinces, -1)
				response.JSON(rw, 0, provinces, "")
				return
			}
		}
	}

	response.JSON(rw, http.StatusBadRequest, nil, "")
}
