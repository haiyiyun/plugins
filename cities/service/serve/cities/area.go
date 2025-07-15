package cities

import (
	"net/http"

	"github.com/haiyiyun/plugins/cities/database/model"
	"github.com/haiyiyun/plugins/cities/database/model/area"
	"github.com/haiyiyun/utils/http/response"
	"go.mongodb.org/mongo-driver/bson"
)

func (self *Service) Route_GET_Area(rw http.ResponseWriter, req *http.Request) {
	cityID := req.URL.Query().Get("city")

	if cityID != "" {
		cacheKey := "cities.area." + cityID
		areas := []model.Area{}
		if found, _ := self.Cache.Get(cacheKey, &areas); found {
			response.JSON(rw, 0, areas, "")
			return
		} else {
			areaModel := area.NewModel(self.M)
			if cur, err := areaModel.Find(req.Context(), bson.D{
				{"city_id", cityID},
			}); err == nil {
				if err := cur.All(req.Context(), &areas); err == nil {
					self.Cache.Set(cacheKey, areas, -1)
					response.JSON(rw, 0, areas, "")
					return
				}
			}
		}
	}

	response.JSON(rw, http.StatusBadRequest, nil, "")
}
