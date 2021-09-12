package cities

import (
	"net/http"

	"github.com/haiyiyun/plugins/cities/database/model"
	"github.com/haiyiyun/plugins/cities/database/model/village"
	"github.com/haiyiyun/utils/http/response"
	"go.mongodb.org/mongo-driver/bson"
)

func (self *Service) Route_GET_Village(rw http.ResponseWriter, req *http.Request) {
	streetID := req.URL.Query().Get("street")

	if streetID != "" {
		cacheKey := "cities.village." + streetID
		if villages, found := self.Cache.Get(cacheKey); found {
			response.JSON(rw, 0, villages, "")
			return
		} else {
			villageModel := village.NewModel(self.M)
			villages := []model.Village{}
			if cur, err := villageModel.Find(req.Context(), bson.D{
				{"street_id", streetID},
			}); err == nil {
				if err := cur.All(req.Context(), &villages); err == nil {
					self.Cache.Set(cacheKey, villages, -1)
					response.JSON(rw, 0, villages, "")
					return
				}
			}
		}

	}

	response.JSON(rw, http.StatusBadRequest, nil, "")
}
