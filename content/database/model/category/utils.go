package category

import (
	"github.com/haiyiyun/mongodb/geometry"
	"github.com/haiyiyun/plugins/content/predefined"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (self *Model) FilterNormalCategory() bson.D {
	filter := bson.D{
		{"enable", true},
		{"status", predefined.PublishStatusNormal},
	}

	return filter
}

func (self *Model) FilterByType(typ int) bson.D {
	filter := bson.D{
		{"type", typ},
	}

	return filter
}

func (self *Model) FilterByParentID(parentID primitive.ObjectID) bson.D {
	filter := bson.D{
		{"parent_id", parentID},
	}

	return filter
}

func (self *Model) FilterByVisibility(visibility int) bson.D {
	filter := bson.D{
		{"visibility", visibility},
	}

	return filter
}

func (self *Model) FilterByVisibilityOrAll(visibility int) bson.D {
	filter := bson.D{
		{"visibility", bson.D{
			{"$in", bson.A{
				visibility,
				predefined.VisibilityTypeAll,
			}},
		}},
	}

	return filter
}

func (self *Model) FilterByTag(tag string) bson.D {
	filter := bson.D{
		{"tags", tag},
	}

	return filter
}

func (self *Model) FilterByLocation(location geometry.Point, maxDistance, minDistance float64) bson.D {
	geo := bson.D{
		{"$geometry", location},
	}

	if maxDistance > 0 {
		geo = append(geo, bson.E{
			"$maxDistance", maxDistance,
		})
	}

	if minDistance > 0 {
		geo = append(geo, bson.E{
			"$minDistance", minDistance,
		})
	}

	filter := bson.D{
		{"$near", geo},
	}

	return filter
}
