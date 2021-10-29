package category

import (
	"net/http"
	"strconv"

	"github.com/haiyiyun/mongodb/geometry"
	"github.com/haiyiyun/plugins/content/database/model/category"
	"github.com/haiyiyun/plugins/content/predefined"
	"github.com/haiyiyun/utils/http/response"
	"github.com/haiyiyun/validator"
	"go.mongodb.org/mongo-driver/bson"
)

func (self *Service) Route_GET_List(rw http.ResponseWriter, r *http.Request) {
	typeStr := r.URL.Query().Get("type")
	typ, _ := strconv.Atoi(typeStr)
	visibilityStr := r.URL.Query().Get("visibility")
	visibility, _ := strconv.Atoi(visibilityStr)
	tags := r.URL.Query()["tags[]"]

	longitudeStr := r.URL.Query().Get("longitude") //经度
	latitudeStr := r.URL.Query().Get("latitude")   //维度
	longitude, _ := strconv.ParseFloat(longitudeStr, 64)
	latitude, _ := strconv.ParseFloat(latitudeStr, 64)
	coordinates := geometry.PointCoordinates{
		longitude, latitude,
	}

	maxDistanceStr := r.URL.Query().Get("max_distance")
	maxDistance, _ := strconv.ParseFloat(maxDistanceStr, 64)
	minDistanceStr := r.URL.Query().Get("min_distance")
	minDistance, _ := strconv.ParseFloat(minDistanceStr, 64)

	valid := validator.Validation{}
	valid.Digital(typeStr).Key("type_str").Message("type必须数字")
	valid.Have(typ,
		predefined.CategoryTypeSystemDynamic,
		predefined.CategoryTypeSystemArticle,
		predefined.CategoryTypeSystemQuestion,
		predefined.CategoryTypeSystemAnswer,
		predefined.CategoryTypeSharePlatform,
	).Key("type").Message("type必须是支持的类型")

	valid.Digital(visibilityStr).Key("visibility_str").Message("visibility必须数字")
	valid.Have(visibility,
		predefined.VisibilityTypeHome,
		predefined.VisibilityTypeRelationship,
		predefined.VisibilityTypeStranger,
		predefined.VisibilityTypeSubject,
		predefined.VisibilityTypeNearly,
		predefined.VisibilityTypeCity,
		predefined.VisibilityTypeProvince,
		predefined.VisibilityTypeNation,
		predefined.VisibilityTypeAll,
	).Key("visibility").Message("visibility必须是支持的类型")

	if valid.HasErrors() {
		response.JSON(rw, http.StatusBadRequest, nil, valid.RandomError().String())
		return
	}

	categoryModel := category.NewModel(self.M)
	filter := categoryModel.FilterNormalCategory()
	filter = append(filter, categoryModel.FilterByType(typ)...)
	filter = append(filter, categoryModel.FilterByVisibilityOrAll(visibility)...)

	if len(tags) > 0 {
		filter = append(filter, bson.E{
			"tags", bson.D{
				{"$in", tags},
			},
		})
	}

	if coordinates != geometry.NilPointCoordinates {
		filter = append(filter, categoryModel.FilterByLocation(geometry.NewPoint(coordinates), maxDistance, minDistance)...)
	}
}