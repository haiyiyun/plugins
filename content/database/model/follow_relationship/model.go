package follow_relationship

import (
	"github.com/haiyiyun/plugins/content/database"
	"github.com/haiyiyun/plugins/content/database/schema"

	"github.com/haiyiyun/mongodb"
)

type Model struct {
	*database.Database `json:"-" bson:"-" map:"-"`
}

func NewModel(mgo mongodb.Mongoer) *Model {
	obj := &Model{
		Database: database.NewDatabase(mgo, schema.FollowRelationship),
	}

	return obj
}
