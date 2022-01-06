package subject

import (
	"github.com/haiyiyun/plugins/content/predefined"

	"github.com/haiyiyun/mongodb/geometry"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (self *Model) FilterNormalSubject() bson.D {
	filter := bson.D{
		{"enable", true},
		{"status", predefined.PublishStatusNormal},
	}

	return filter
}

func (self *Model) FilterByPublishUserID(publishUserID primitive.ObjectID) bson.D {
	filter := bson.D{
		{"publish_user_id", publishUserID},
	}

	return filter
}

func (self *Model) FilterByPublishUserIDs(publishUserIDs []primitive.ObjectID) bson.D {
	if len(publishUserIDs) == 0 {
		return bson.D{}
	}

	filter := bson.D{
		{"publish_user_id", bson.D{
			{"$in", publishUserIDs},
		}},
	}

	return filter
}

func (self *Model) FilterByExcludePublishUserIDs(publishUserIDs []primitive.ObjectID) bson.D {
	if len(publishUserIDs) == 0 {
		return bson.D{}
	}

	filter := bson.D{
		{"publish_user_id", bson.D{
			{"$nin", publishUserIDs},
		}},
	}

	return filter
}

func (self *Model) FilterByType(typ int) bson.D {
	filter := bson.D{
		{"type", typ},
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
		{"location", bson.D{
			{"$near", geo},
		}},
	}

	return filter
}
