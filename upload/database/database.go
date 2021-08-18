package database

import (
	"context"
	"time"

	"github.com/haiyiyun/mongodb"
	"github.com/haiyiyun/mongodb/driver"
	"github.com/haiyiyun/utils/help"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	*driver.Driver `json:"-" bson:"-" map:"-"`
}

func NewDatabase(mgo mongodb.Mongoer, col bson.M) *Database {
	mdl := &Database{
		Driver: driver.NewDriver(mgo, col),
	}

	return mdl
}

func (mdl *Database) UpdateOne(ctx context.Context, filter, update bson.D, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	update = append(update, mdl.DataSet(bson.D{
		{"update_time", time.Now()},
	})...)

	return mdl.Driver.UpdateOne(ctx, filter, update, opts...)
}

func (mdl *Database) Set(ctx context.Context, filter, data bson.D, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return mdl.UpdateOne(ctx, filter, mdl.DataSet(data), opts...)
}

func (mdl *Database) SetOnInsert(ctx context.Context, filter, update bson.D, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	insertData, opt := mdl.DataSetOnInsert(bson.D{
		{"create_time", time.Now()},
	}, opts...)

	return mdl.UpdateOne(ctx, filter, append(update, insertData...), opt)
}

func (mdl *Database) AddToSet(ctx context.Context, filter, addToSetData bson.D, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return mdl.UpdateOne(ctx, filter, mdl.DataAddToSet(addToSetData), opts...)
}

func (mdl *Database) AddToSetOnInsert(ctx context.Context, filter, addToSetData bson.D, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return mdl.SetOnInsert(ctx, filter, mdl.DataAddToSet(addToSetData), opts...)
}

func (mdl *Database) SetAndAddToSet(ctx context.Context, filter, setData, addToSetData bson.D, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return mdl.UpdateOne(ctx, filter, append(mdl.DataSet(setData), mdl.DataAddToSet(addToSetData)...), opts...)
}

func (mdl *Database) SetAndAddToSetOnInsert(ctx context.Context, filter, setData, addToSetData bson.D, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return mdl.SetOnInsert(ctx, filter, append(mdl.DataSet(setData), mdl.DataAddToSet(addToSetData)...), opts...)
}

func (mdl *Database) Pull(ctx context.Context, filter, pullFilters bson.D, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return mdl.UpdateOne(ctx, filter, mdl.DataPull(pullFilters), opts...)
}

func (mdl *Database) Create(ctx context.Context, m interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	mm := help.NewStruct(m).StructToMap()

	if id, found := mm["_id"]; found {
		if id == primitive.NilObjectID {
			mm["_id"] = primitive.NewObjectID()
		}
	}

	if ct, found := mm["create_time"].(time.Time); found {
		if ct.IsZero() {
			mm["create_time"] = time.Now()
		}
	}

	if ut, found := mm["update_time"].(time.Time); found {
		if ut.IsZero() {
			mm["update_time"] = time.Now()
		}
	}

	return mdl.InsertOne(ctx, mm, opts...)
}
