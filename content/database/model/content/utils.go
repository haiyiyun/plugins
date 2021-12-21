package content

import (
	"context"
	"time"

	"github.com/haiyiyun/mongodb/geometry"
	"github.com/haiyiyun/plugins/content/predefined"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (self *Model) FilterNormalContent() bson.D {
	filter := bson.D{
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

func (self *Model) FilterByCategoryID(categoryID primitive.ObjectID) bson.D {
	filter := bson.D{
		{"category_id", categoryID},
	}

	return filter
}

func (self *Model) FilterBySubjectID(subjectID primitive.ObjectID) bson.D {
	filter := bson.D{
		{"subject_id", subjectID},
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

func (self *Model) FilterByPublishType(publishType int) bson.D {
	filter := bson.D{
		{"publish_type", publishType},
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

func (self *Model) FilterByDiscussEstimateTotal(discussEstimateTotal int) bson.D {
	filter := bson.D{
		{"discuss_estimate_total", discussEstimateTotal},
	}

	return filter
}

func (self *Model) FilterByGteDiscussEstimateTotal(discussEstimateTotal int) bson.D {
	filter := bson.D{
		{"discuss_estimate_total", bson.D{
			{"$gte", discussEstimateTotal},
		}},
	}

	return filter
}

func (self *Model) FilterByLteDiscussEstimateTotal(discussEstimateTotal int) bson.D {
	filter := bson.D{
		{"discuss_estimate_total", bson.D{
			{"$lte", discussEstimateTotal},
		}},
	}

	return filter
}

func (self *Model) FilterByValue(value int) bson.D {
	filter := bson.D{
		{"value", value},
	}

	return filter
}

func (self *Model) FilterByGteValue(value int) bson.D {
	filter := bson.D{
		{"value", bson.D{
			{"$gte", value},
		}},
	}

	return filter
}

func (self *Model) FilterByLteValue(value int) bson.D {
	filter := bson.D{
		{"value", bson.D{
			{"$lte", value},
		}},
	}

	return filter
}

func (self *Model) FilterByGteStartTime(startTime time.Time) bson.D {
	filter := bson.D{
		{"start_time", bson.D{
			{"$gte", startTime},
		}},
	}

	return filter
}

func (self *Model) FilterByLteEndTime(endTime time.Time) bson.D {
	filter := bson.D{
		{"end_time", bson.D{
			{"$lte", endTime},
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

func (self *Model) AddAtUsers(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.AddToSet(cxt, filter, bson.D{
		{"at_users", userID},
	})
}

func (self *Model) DeleteAtUsers(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.Pull(cxt, filter, bson.D{
		{"at_users", userID},
	})
}

func (self *Model) AddOnlyUserIDShowDetail(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.AddToSet(cxt, filter, bson.D{
		{"only_user_id_show_detail", userID},
	})
}

func (self *Model) DeleteOnlyUserIDShowDetail(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.Pull(cxt, filter, bson.D{
		{"only_user_id_show_detail", userID},
	})
}

func (self *Model) AddOnlyUserIDDiscuss(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.AddToSet(cxt, filter, bson.D{
		{"only_user_id_discuss", userID},
	})
}

func (self *Model) DeleteOnlyUserIDDiscuss(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.Pull(cxt, filter, bson.D{
		{"only_user_id_discuss", userID},
	})
}

func (self *Model) AddOnlyUserIDCanReplyDiscuss(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.AddToSet(cxt, filter, bson.D{
		{"only_publish_user_id_can_reply_discuss", userID},
	})
}

func (self *Model) DeleteOnlyUserIDCanReplyDiscuss(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.Pull(cxt, filter, bson.D{
		{"only_publish_user_id_can_reply_discuss", userID},
	})
}

func (self *Model) AddOnlyUserIDCanNotReplyDiscuss(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.AddToSet(cxt, filter, bson.D{
		{"only_publish_user_id_can_not_reply_discuss", userID},
	})
}

func (self *Model) DeleteOnlyUserIDCanNotReplyDiscuss(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.Pull(cxt, filter, bson.D{
		{"only_publish_user_id_can_not_reply_discuss", userID},
	})
}

func (self *Model) AddOnlyUserIDShowDiscuss(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.AddToSet(cxt, filter, bson.D{
		{"only_user_id_show_discuss", userID},
	})
}

func (self *Model) DeleteOnlyUserIDShowDiscuss(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.Pull(cxt, filter, bson.D{
		{"only_user_id_show_discuss", userID},
	})
}

func (self *Model) AddReadedUser(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.AddToSet(cxt, filter, bson.D{
		{"readed_user", userID},
	})
}

func (self *Model) DeleteReadedUser(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.Pull(cxt, filter, bson.D{
		{"readed_user", userID},
	})
}

func (self *Model) AddWantedUser(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.AddToSet(cxt, filter, bson.D{
		{"wanted_user", userID},
	})
}

func (self *Model) DeleteWantedUser(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.Pull(cxt, filter, bson.D{
		{"wanted_user", userID},
	})
}

func (self *Model) AddLikedUser(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.AddToSet(cxt, filter, bson.D{
		{"liked_user", userID},
	})
}

func (self *Model) DeleteLikedUser(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.Pull(cxt, filter, bson.D{
		{"liked_user", userID},
	})
}

func (self *Model) AddHatedUser(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.AddToSet(cxt, filter, bson.D{
		{"hated_user", userID},
	})
}

func (self *Model) DeleteHatedUser(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.Pull(cxt, filter, bson.D{
		{"hated_user", userID},
	})
}

func (self *Model) AddAntiGuiseUser(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.AddToSet(cxt, filter, bson.D{
		{"anti_guise_user", userID},
	})
}

func (self *Model) DeleteAntiGuiseUser(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.Pull(cxt, filter, bson.D{
		{"anti_guise_user", userID},
	})
}
