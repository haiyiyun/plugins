package log

import (
	"context"
	"net/http"

	"github.com/haiyiyun/plugins/log/database/model"
	modelLog "github.com/haiyiyun/plugins/log/database/model/log"

	"github.com/haiyiyun/utils/http/pagination"
	"github.com/haiyiyun/utils/http/response"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (self *Service) Route_GET_Index(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	logModel := modelLog.NewModel(self.M)
	filter := bson.D{}

	types := r.Form["type[]"]
	if len(types) > 0 {
		filter = append(filter, bson.E{
			"type", bson.D{
				{"$in", types},
			},
		})
	}

	methods := r.Form["method[]"]
	if len(methods) > 0 {
		filter = append(filter, bson.E{
			"method", bson.D{
				{"$in", methods},
			},
		})
	}

	if userIDHex := r.FormValue("user_id"); userIDHex != "" {
		if userID, err := primitive.ObjectIDFromHex(userIDHex); err == nil {
			filter = append(filter, bson.E{"user_id", userID})
		}
	}

	if user := r.FormValue("user"); user != "" {
		filter = append(filter, bson.E{"user", user})
	}

	if ip := r.FormValue("ip"); ip != "" {
		filter = append(filter, bson.E{"ip", ip})
	}

	if path := r.FormValue("path"); path != "" {
		filter = append(filter, bson.E{"path", path})
	}

	cnt, _ := logModel.CountDocuments(context.Background(), filter)
	pg := pagination.Parse(r, cnt)

	opt := options.Find().SetSort(bson.D{
		{"create_time", -1},
	}).SetProjection(bson.D{}).SetSkip(pg.SkipNum).SetLimit(pg.PageSize)

	items := []model.Log{}
	if cur, err := logModel.Find(context.Background(), filter, opt); err == nil {
		cur.All(context.TODO(), &items)
	}

	rpr := response.ResponsePaginationResult{
		Total: cnt,
		Items: items,
	}

	response.JSON(rw, 0, rpr, "")
}
