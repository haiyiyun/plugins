package category

import (
	"net/http"

	"github.com/haiyiyun/log"
	"github.com/haiyiyun/mongodb/geometry"
	"github.com/haiyiyun/plugins/content/database/model/category"
	"github.com/haiyiyun/plugins/content/predefined"
	"github.com/haiyiyun/utils/help"
	"github.com/haiyiyun/utils/http/pagination"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/utils/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (self *Service) Route_GET_List(rw http.ResponseWriter, r *http.Request) {
	var requestCL predefined.RequestServeCategoryList
	if err := validator.FormStruct(&requestCL, r.URL.Query()); err != nil {
		response.JSON(rw, http.StatusBadRequest, nil, err.Error())
		return
	}

	coordinates := geometry.PointCoordinates{
		requestCL.Longitude, requestCL.Latitude,
	}

	categoryModel := category.NewModel(self.M)
	filter := categoryModel.FilterNormalCategory()
	filter = append(filter, categoryModel.FilterByType(requestCL.Type)...)
	filter = append(filter, categoryModel.FilterByVisibility(requestCL.Visibility)...)

	if requestCL.ParentID != primitive.NilObjectID {
		filter = append(filter, categoryModel.FilterByParentID(requestCL.ParentID)...)
	}

	if len(requestCL.Tags) > 0 {
		filter = append(filter, bson.E{
			"tags", bson.D{
				{"$in", requestCL.Tags},
			},
		})
	}

	if coordinates != geometry.NilPointCoordinates {
		filter = append(filter, categoryModel.FilterByLocation(geometry.NewPoint(coordinates), requestCL.MaxDistance, requestCL.MinDistance)...)
	}

	cnt, _ := categoryModel.CountDocuments(r.Context(), filter)
	pg := pagination.Parse(r, cnt)

	projection := bson.D{}

	opt := options.Find().SetSort(bson.D{
		{"create_time", -1},
	}).SetProjection(projection).SetSkip(pg.SkipNum).SetLimit(pg.PageSize)

	if cur, err := categoryModel.Find(r.Context(), filter, opt); err != nil {
		response.JSON(rw, http.StatusNotFound, nil, "")
	} else {
		items := []help.M{}
		if err := cur.All(r.Context(), &items); err != nil {
			log.Error(err)
			response.JSON(rw, http.StatusServiceUnavailable, nil, "")
		} else {
			rpr := response.ResponsePaginationResult{
				Total: cnt,
				Items: items,
			}

			response.JSON(rw, 0, rpr, "")
		}
	}
}
