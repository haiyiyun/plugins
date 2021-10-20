package profile_audit

import (
	"github.com/haiyiyun/plugins/user_profile/database"
	"github.com/haiyiyun/plugins/user_profile/database/schema"

	"github.com/haiyiyun/mongodb"
)

type Model struct {
	*database.Database `json:"-" bson:"-" map:"-"`
}

func NewModel(mgo mongodb.Mongoer) *Model {
	obj := &Model{
		Database: database.NewDatabase(mgo, schema.ProfileAudit),
	}

	return obj
}
