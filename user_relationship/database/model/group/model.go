package group

import (
	"github.com/haiyiyun/plugins/user_relationship/database"
	"github.com/haiyiyun/plugins/user_relationship/database/schema"

	"github.com/haiyiyun/mongodb"
)

type Model struct {
	*database.Database `json:"-" bson:"-" map:"-"`
}

func NewModel(mgo mongodb.Mongoer) *Model {
	obj := &Model{
		Database: database.NewDatabase(mgo, schema.Group),
	}

	return obj
}
