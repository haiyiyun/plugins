package token

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
		Database: database.NewDatabase(mgo, schema.Token),
	}

	return obj
}
