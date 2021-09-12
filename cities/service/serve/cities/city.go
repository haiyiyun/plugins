package cities

import (
	"net/http"

	"github.com/haiyiyun/plugins/cities/database/model"
	"github.com/haiyiyun/plugins/cities/database/model/city"
	"github.com/haiyiyun/utils/http/response"
	"go.mongodb.org/mongo-driver/bson"
)

func (self *Service) Route_GET_City(rw http.ResponseWriter, req *http.Request) {
	provinceID := req.URL.Query().Get("province")

	if provinceID != "" {
		cityModel := city.NewModel(self.M)
		cities := []model.City{}
		if cur, err := cityModel.Find(req.Context(), bson.D{
			{"province_id", provinceID},
		}); err == nil {
			if err := cur.All(req.Context(), &cities); err == nil {
				response.JSON(rw, 0, cities, "")
				return
			}
		}
	}

	response.JSON(rw, http.StatusBadRequest, nil, "")
}
