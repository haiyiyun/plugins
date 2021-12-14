package predefined

import "go.mongodb.org/mongo-driver/bson/primitive"

type RequestServeVisitor struct {
	OwnerUserID primitive.ObjectID `form:"owner_user_id" validate:"required"`
}
