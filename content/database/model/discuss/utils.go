package discuss

import (
	"context"

	"github.com/haiyiyun/mongodb/geometry"
	"github.com/haiyiyun/plugins/content/predefined"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (self *Model) FilterNormalDiscuss() bson.D {
	filter := bson.D{
		{"status", predefined.PublishStatusNormal},
	}

	return filter
}

func (self *Model) FilterByObjectID(objectID primitive.ObjectID) bson.D {
	filter := bson.D{
		{"object_id", objectID},
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

func (self *Model) FilterByType(typ int) bson.D {
	filter := bson.D{
		{"type", typ},
	}

	return filter
}

func (self *Model) FilterByTypes(types []int) bson.D {
	if len(types) == 0 {
		return bson.D{}
	}

	filter := bson.D{
		{"type", bson.D{
			{"$in", types},
		}},
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

func (self *Model) AddLikedUser(cxt context.Context, discussID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(discussID)
	return self.AddToSet(cxt, filter, bson.D{
		{"liked_user", userID},
	})
}

func (self *Model) DeleteLikedUser(cxt context.Context, discussID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(discussID)
	return self.Pull(cxt, filter, bson.D{
		{"liked_user", userID},
	})
}

func (self *Model) AddHatedUser(cxt context.Context, discussID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(discussID)
	return self.AddToSet(cxt, filter, bson.D{
		{"hated_user", userID},
	})
}

func (self *Model) DeleteHatedUser(cxt context.Context, discussID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(discussID)
	return self.Pull(cxt, filter, bson.D{
		{"hated_user", userID},
	})
}

func (self *Model) FilterByEvaluation(evaluation int) bson.D {
	filter := bson.D{
		{"evaluation", evaluation},
	}

	return filter
}

func (self *Model) FilterByGteEvaluation(evaluation int) bson.D {
	filter := bson.D{
		{"evaluation", bson.D{
			{"$gte", evaluation},
		}},
	}

	return filter
}

func (self *Model) FilterByLteEvaluation(evaluation int) bson.D {
	filter := bson.D{
		{"evaluation", bson.D{
			{"$lte", evaluation},
		}},
	}

	return filter
}
