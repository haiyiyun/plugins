package dictionary

import (
	"github.com/haiyiyun/plugins/dictionary/database"
	"github.com/haiyiyun/plugins/dictionary/database/schema"

	"github.com/haiyiyun/mongodb"
)

type Model struct {
	*database.Database `json:"-" bson:"-" map:"-"`
}

func NewModel(mgo mongodb.Mongoer) *Model {
	obj := &Model{
		Database: database.NewDatabase(mgo, schema.Dictionary),
	}

	return obj
}
