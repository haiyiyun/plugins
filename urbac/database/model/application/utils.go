package application

import (
	"encoding/gob"

	"github.com/haiyiyun/plugins/urbac/database/model"

	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	gob.Register(model.Application{})
	gob.Register([]model.Application{})
	gob.Register(map[string]model.Application{})
	gob.Register(model.ApplicationModule{})
	gob.Register(model.ApplicationModuleAction{})
}

func (self *Model) FilterByName(name string) bson.D {
	return bson.D{
		{"name", name},
	}
}
