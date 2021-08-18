package log

import (
	"github.com/haiyiyun/plugins/log/database"
	"github.com/haiyiyun/plugins/log/database/schema"

	"github.com/haiyiyun/mongodb"
)

type Model struct {
	*database.Database `json:"-" bson:"-" map:"-"`
}

func NewModel(mgo mongodb.Mongoer) *Model {
	obj := &Model{
		Database: database.NewDatabase(mgo, schema.Log),
	}

	return obj
}
