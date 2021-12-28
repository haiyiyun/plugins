package follow_content

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (self *Model) FilterByFollowRelationshipID(followRelationshipID primitive.ObjectID) bson.D {
	return bson.D{
		{"follow_relationship_id", followRelationshipID},
	}
}

func (self *Model) FilterByUserID(userID primitive.ObjectID) bson.D {
	return bson.D{
		{"user_id", userID},
	}
}

func (self *Model) FilterByContentID(contentID primitive.ObjectID) bson.D {
	return bson.D{
		{"content_id", contentID},
	}
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

func (self *Model) FilterByType(typ int) bson.D {
	return bson.D{
		{"type", typ},
	}
}

func (self *Model) FilterByExtensionID(extensionID primitive.ObjectID) bson.D {
	return bson.D{
		{"extension_id", extensionID},
	}
}

func (self *Model) SetReadedTimeByID(ctx context.Context, id primitive.ObjectID, readedTime time.Time) (*mongo.UpdateResult, error) {
	return self.Set(ctx, self.FilterByID(id), bson.D{
		{"readed_time", readedTime},
	})
}

func (self *Model) SetAllReadedTimeByUserIDAndContentID(ctx context.Context, userID, contentID primitive.ObjectID, readedTime time.Time) (*mongo.UpdateResult, error) {
	filter := self.FilterByUserID(userID)
	filter = append(filter, self.FilterByContentID(contentID)...)
	return self.UpdateMany(ctx, filter, self.DataSet(bson.D{
		{"readed_time", readedTime},
		{"update_time", time.Now()},
	}))
}
