package country

import (
	"github.com/haiyiyun/plugins/cities/database"
	"github.com/haiyiyun/plugins/cities/database/schema"

	"github.com/haiyiyun/mongodb"
)

type Model struct {
	*database.Database `json:"-" bson:"-" map:"-"`
}

func NewModel(mgo mongodb.Mongoer) *Model {
	obj := &Model{
		Database: database.NewDatabase(mgo, schema.Country),
	}

	return obj
}
