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

func (self *Model) FilterByAssociateType(associateType int) bson.D {
	filter := bson.D{
		{"associate_type", associateType},
	}

	return filter
}

func (self *Model) FilterByAssociateID(associateID primitive.ObjectID) bson.D {
	filter := bson.D{
		{"associate_id", associateID},
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

func (self *Model) FilterByLimitUserRole(limitUserRole []string) bson.D {
	if len(limitUserRole) == 0 {
		return bson.D{}
	}

	filter := bson.D{
		{"limit_user_role", bson.D{
			{"$in", limitUserRole},
		}},
	}

	return filter
}

func (self *Model) FilterByLimitUserTag(limitUserTag []string) bson.D {
	if len(limitUserTag) == 0 {
		return bson.D{}
	}

	filter := bson.D{
		{"limit_user_tag", bson.D{
			{"$in", limitUserTag},
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

func (self *Model) FilterByGteLimitUserAtLeastLevel(limitUserAtLeastLevel int) bson.D {
	filter := bson.D{
		{"limit_user_at_least_level", bson.D{
			{"$gte", limitUserAtLeastLevel},
		}},
	}

	return filter
}

func (self *Model) FilterByStartTimeGte(startTime time.Time) bson.D {
	filter := bson.D{
		{"start_time", bson.D{
			{"$gte", startTime},
		}},
	}

	return filter
}

func (self *Model) FilterByStartTimeLte(startTime time.Time) bson.D {
	filter := bson.D{
		{"start_time", bson.D{
			{"$lte", startTime},
		}},
	}

	return filter
}

func (self *Model) FilterByEndTimeGte(endTime time.Time) bson.D {
	filter := bson.D{
		{"end_time", bson.D{
			{"$gte", endTime},
		}},
	}

	return filter
}

func (self *Model) FilterByEndTimeLte(endTime time.Time) bson.D {
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
		{"location", bson.D{
			{"$near", geo},
		}},
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
	filter = append(filter, self.FilterNormalContent()...)
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

func (self *Model) AddOnlyUserIDNotLimitUserLevel(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.AddToSet(cxt, filter, bson.D{
		{"only_user_id_not_limit_user_level", userID},
	})
}

func (self *Model) DeleteOnlyUserIDNotLimitUserLevel(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.Pull(cxt, filter, bson.D{
		{"only_user_id_not_limit_user_level", userID},
	})
}

func (self *Model) AddOnlyUserIDNotLimitUserRole(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.AddToSet(cxt, filter, bson.D{
		{"only_user_id_not_limit_user_role", userID},
	})
}

func (self *Model) DeleteOnlyUserIDNotLimitUserRole(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.Pull(cxt, filter, bson.D{
		{"only_user_id_not_limit_user_role", userID},
	})
}

func (self *Model) AddOnlyUserIDNotLimitUserTag(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.AddToSet(cxt, filter, bson.D{
		{"only_user_id_not_limit_user_tag", userID},
	})
}

func (self *Model) DeleteOnlyUserIDNotLimitUserTag(cxt context.Context, contentID, userID primitive.ObjectID) (*mongo.UpdateResult, error) {
	filter := self.FilterByID(contentID)
	return self.Pull(cxt, filter, bson.D{
		{"only_user_id_not_limit_user_tag", userID},
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
