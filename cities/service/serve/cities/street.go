package cities

import (
	"net/http"

	"github.com/haiyiyun/plugins/cities/database/model"
	"github.com/haiyiyun/plugins/cities/database/model/street"
	"github.com/haiyiyun/utils/http/response"
	"go.mongodb.org/mongo-driver/bson"
)

func (self *Service) Route_GET_Street(rw http.ResponseWriter, req *http.Request) {
	areaID := req.URL.Query().Get("area")

	if areaID != "" {
		streetModel := street.NewModel(self.M)
		streets := []model.Street{}
		if cur, err := streetModel.Find(req.Context(), bson.D{
			{"area_id", areaID},
		}); err == nil {
			if err := cur.All(req.Context(), &streets); err == nil {
				response.JSON(rw, 0, streets, "")
				return
			}
		}
	}

	response.JSON(rw, http.StatusBadRequest, nil, "")
}
