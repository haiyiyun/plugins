package role

import (
	"github.com/haiyiyun/plugins/urbac/database"
	"github.com/haiyiyun/plugins/urbac/database/schema"

	"github.com/haiyiyun/mongodb"
)

type Model struct {
	*database.Database `json:"-" bson:"-" map:"-"`
}

func NewModel(mgo mongodb.Mongoer) *Model {
	obj := &Model{
		Database: database.NewDatabase(mgo, schema.Role),
	}

	return obj
}
