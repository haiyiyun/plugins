package follow_relationship

import (
	"context"

	"github.com/haiyiyun/plugins/content/database/model"
	"github.com/haiyiyun/plugins/content/predefined"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (self *Model) FilterByRelationship(typ int, userID, objectID primitive.ObjectID) bson.D {
	return bson.D{
		{"type", typ},
		{"user_id", userID},
		{"object_id", objectID},
	}
}

func (self *Model) FilterByUserWithType(userID primitive.ObjectID, typ int) bson.D {
	return bson.D{
		{"user_id", userID},
		{"type", typ},
	}
}

func (self *Model) FilterByObjectIDWithType(objectID primitive.ObjectID, typ int) bson.D {
	return bson.D{
		{"object_id", objectID},
		{"type", typ},
	}
}

func (self *Model) CreateRelationship(ctx context.Context, typ int, userID, objectID primitive.ObjectID, stealth bool) (primitive.ObjectID, error) {
	var id primitive.ObjectID
	err := self.UseSession(ctx, func(sctx mongo.SessionContext) error {
		if err := sctx.StartTransaction(); err != nil {
			return err
		}

		mutual := false

		if typ == predefined.FollowTypeUser {
			if sr := self.FindOne(sctx, self.FilterByRelationship(typ, objectID, userID), options.FindOne().SetProjection(bson.D{
				{"_id", 1},
			})); sr.Err() != nil {
				if sr.Err() != mongo.ErrNoDocuments {
					sctx.AbortTransaction(sctx)
					return sr.Err()
				}
			} else {
				var frel model.FollowRelationship
				if err := sr.Decode(&frel); err != nil {
					sctx.AbortTransaction(sctx)
					return err
				} else {
					if _, err := self.Set(sctx, self.FilterByID(frel.ID), bson.D{
						{"mutual", true},
					}); err == nil {
						mutual = true
					} else {
						sctx.AbortTransaction(sctx)
						return err
					}
				}
			}
		}

		rel := &model.FollowRelationship{
			Type:     typ,
			UserID:   userID,
			ObjectID: objectID,
			Mutual:   mutual,
			Stealth:  stealth,
		}

		ior, err := self.Create(sctx, rel)
		if err != nil {
			sctx.AbortTransaction(sctx)
			return err
		}

		id = ior.InsertedID.(primitive.ObjectID)

		return sctx.CommitTransaction(sctx)
	})

	return id, err
}

func (self *Model) DeleteRelationship(ctx context.Context, typ int, userID, objectID primitive.ObjectID) error {
	err := self.UseSession(ctx, func(sctx mongo.SessionContext) error {
		if typ == predefined.FollowTypeUser {
			if sr := self.FindOne(sctx, self.FilterByRelationship(typ, objectID, userID), options.FindOne().SetProjection(bson.D{
				{"_id", 1},
			})); sr.Err() != nil {
				if sr.Err() != mongo.ErrNoDocuments {
					sctx.AbortTransaction(sctx)
					return sr.Err()
				}
			} else {
				var frel model.FollowRelationship
				if err := sr.Decode(&frel); err != nil {
					sctx.AbortTransaction(sctx)
					return err
				} else {
					if _, err := self.Set(sctx, self.FilterByID(frel.ID), bson.D{
						{"mutual", false},
					}); err != nil {
						sctx.AbortTransaction(sctx)
						return err
					}
				}
			}
		}

		_, err := self.DeleteOne(ctx, self.FilterByRelationship(typ, userID, objectID))
		if err != nil {
			sctx.AbortTransaction(sctx)
			return err
		}

		return sctx.CommitTransaction(sctx)
	})

	return err
}

func (self *Model) DeleteRelationshipByID(ctx context.Context, id primitive.ObjectID) error {
	if sr := self.FindOne(ctx, self.FilterByID(id)); sr.Err() != nil {
		return sr.Err()
	} else {
		var frel model.FollowRelationship
		if err := sr.Decode(&frel); err != nil {
			return err
		} else {
			return self.DeleteRelationship(ctx, frel.Type, frel.UserID, frel.ObjectID)
		}
	}
}
