package upload

import (
	"github.com/haiyiyun/plugins/upload/database"
	"github.com/haiyiyun/plugins/upload/database/schema"

	"github.com/haiyiyun/mongodb"
)

type Model struct {
	*database.Database `json:"-" bson:"-" map:"-"`
}

func NewModel(mgo mongodb.Mongoer) *Model {
	obj := &Model{
		Database: database.NewDatabase(mgo, schema.Upload),
	}

	return obj
}
