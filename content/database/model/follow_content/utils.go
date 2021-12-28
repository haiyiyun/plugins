package follow_content

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
