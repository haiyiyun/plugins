package contacts_blacklist

import (
	"github.com/haiyiyun/plugins/user/database"
	"github.com/haiyiyun/plugins/user/database/schema"

	"github.com/haiyiyun/mongodb"
)

type Model struct {
	*database.Database `json:"-" bson:"-" map:"-"`
}

func NewModel(mgo mongodb.Mongoer) *Model {
	obj := &Model{
		Database: database.NewDatabase(mgo, schema.ContactsBlacklist),
	}

	return obj
}
