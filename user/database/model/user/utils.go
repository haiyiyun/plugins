package user

import (
	"context"
	"encoding/gob"
	"time"

	"github.com/haiyiyun/plugins/user/database/model"
	"github.com/haiyiyun/utils/help"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func init() {
	gob.Register(model.User{})
}

func (self *Model) FilterByName(name string) bson.D {
	return bson.D{
		{"name", name},
	}
}

func (self *Model) CheckNameAndPassword(name, password string) (u model.User, err error) {
	passwordMd5 := help.Strings(password).Md5()

	filter := bson.D{
		{"name", name},
		{"password", passwordMd5},
		{"enable", true},
		{"delete", false},
	}

	sr := self.FindOne(context.TODO(), filter)
	err = sr.Decode(&u)

	return
}

func (self *Model) GetUserByID(userID primitive.ObjectID) (u model.User, err error) {
	filter := self.FilterByID(userID)
	filter = append(filter, bson.D{
		{"enable", true},
		{"delete", false},
	}...)
	sr := self.FindOne(context.TODO(), filter)
	err = sr.Decode(&u)

	return
}

func (self *Model) ChangePassword(userID primitive.ObjectID, password string) error {
	_, err := self.Set(context.TODO(), self.FilterByID(userID), bson.D{
		{"password", help.NewString(password).Md5()},
	})

	return err
}

func (self *Model) GuestToUser(userID primitive.ObjectID, username, password string) error {
	_, err := self.Set(context.TODO(), self.FilterByID(userID), bson.D{
		{"name", username},
		{"password", help.NewString(password).Md5()},
		{"guest", false},
	})

	return err
}

func (self *Model) AddRole(cxt context.Context, userID primitive.ObjectID, role string, startTime, endTime time.Time) error {
	roleData := model.UserRole{
		Role:      role,
		StartTime: startTime,
		EndTime:   endTime,
	}

	_, err := self.AddToSet(cxt, self.FilterByID(userID), bson.D{
		{"roles", roleData},
	})

	return err
}

func (self *Model) UpdateRoleStartTime(cxt context.Context, userID primitive.ObjectID, role string, startTime time.Time) error {
	filter := self.FilterByID(userID)
	filter = append(filter, bson.D{
		{"roles.role", role},
	}...)

	_, err := self.Set(cxt, filter, bson.D{
		{"roles.$.start_time", startTime},
	})

	return err
}

func (self *Model) UpdateRoleEndTime(cxt context.Context, userID primitive.ObjectID, role string, endTime time.Time) error {
	filter := self.FilterByID(userID)
	filter = append(filter, bson.D{
		{"roles.role", role},
	}...)

	_, err := self.Set(cxt, filter, bson.D{
		{"roles.$.end_time", endTime},
	})

	return err
}

func (self *Model) UpdateRoleLevel(cxt context.Context, userID primitive.ObjectID, role string, level int) error {
	filter := self.FilterByID(userID)
	filter = append(filter, bson.D{
		{"roles.role", role},
	}...)

	_, err := self.Set(cxt, filter, bson.D{
		{"roles.$.level", level},
	})

	return err
}

func (self *Model) UpdateRoleIncLevel(cxt context.Context, userID primitive.ObjectID, role string, incLevel int) error {
	filter := self.FilterByID(userID)
	filter = append(filter, bson.D{
		{"roles.role", role},
	}...)

	_, err := self.Inc(cxt, filter, bson.D{
		{"roles.$.level", incLevel},
	})

	return err
}

func (self *Model) UpdateRoleExperience(cxt context.Context, userID primitive.ObjectID, role string, experience int) error {
	filter := self.FilterByID(userID)
	filter = append(filter, bson.D{
		{"roles.role", role},
	}...)

	_, err := self.Set(cxt, filter, bson.D{
		{"roles.$.experience", experience},
	})

	return err
}

func (self *Model) UpdateRoleIncExperience(cxt context.Context, userID primitive.ObjectID, role string, incExperience int) error {
	filter := self.FilterByID(userID)
	filter = append(filter, bson.D{
		{"roles.role", role},
	}...)

	_, err := self.Inc(cxt, filter, bson.D{
		{"roles.$.experience", incExperience},
	})

	return err
}

func (self *Model) DeleteRole(cxt context.Context, userID primitive.ObjectID, role string) error {
	_, err := self.Pull(cxt, self.FilterByID(userID), bson.D{
		{"roles", bson.D{
			{"$elemMatch", bson.D{
				{"role", role},
			}},
		}},
	})

	return err
}

func (self *Model) AddTag(cxt context.Context, userID primitive.ObjectID, tag string, startTime, endTime time.Time) error {
	tagData := model.UserTag{
		Tag:       tag,
		StartTime: startTime,
		EndTime:   endTime,
	}

	_, err := self.AddToSet(cxt, self.FilterByID(userID), bson.D{
		{"tags", tagData},
	})

	return err
}

func (self *Model) UpdateTagStartTime(cxt context.Context, userID primitive.ObjectID, tag string, startTime time.Time) error {
	filter := self.FilterByID(userID)
	filter = append(filter, bson.D{
		{"tags.tag", tag},
	}...)

	_, err := self.Set(cxt, filter, bson.D{
		{"tags.$.start_time", startTime},
	})

	return err
}

func (self *Model) UpdateTagEndTime(cxt context.Context, userID primitive.ObjectID, tag string, endTime time.Time) error {
	filter := self.FilterByID(userID)
	filter = append(filter, bson.D{
		{"tags.tag", tag},
	}...)

	_, err := self.Set(cxt, filter, bson.D{
		{"tags.$.end_time", endTime},
	})

	return err
}

func (self *Model) UpdateTagLevel(cxt context.Context, userID primitive.ObjectID, tag string, level int) error {
	filter := self.FilterByID(userID)
	filter = append(filter, bson.D{
		{"tags.tag", tag},
	}...)

	_, err := self.Set(cxt, filter, bson.D{
		{"tags.$.level", level},
	})

	return err
}

func (self *Model) UpdateTagIncLevel(cxt context.Context, userID primitive.ObjectID, tag string, incLevel int) error {
	filter := self.FilterByID(userID)
	filter = append(filter, bson.D{
		{"tags.tag", tag},
	}...)

	_, err := self.Inc(cxt, filter, bson.D{
		{"tags.$.level", incLevel},
	})

	return err
}

func (self *Model) UpdateTagExperience(cxt context.Context, userID primitive.ObjectID, tag string, experience int) error {
	filter := self.FilterByID(userID)
	filter = append(filter, bson.D{
		{"tags.tag", tag},
	}...)

	_, err := self.Set(cxt, filter, bson.D{
		{"tags.$.experience", experience},
	})

	return err
}

func (self *Model) UpdateTagIncExperience(cxt context.Context, userID primitive.ObjectID, tag string, incExperience int) error {
	filter := self.FilterByID(userID)
	filter = append(filter, bson.D{
		{"tags.tag", tag},
	}...)

	_, err := self.Inc(cxt, filter, bson.D{
		{"tags.$.experience", incExperience},
	})

	return err
}

func (self *Model) DeleteTag(cxt context.Context, userID primitive.ObjectID, tag string) error {
	_, err := self.Pull(cxt, self.FilterByID(userID), bson.D{
		{"tags", bson.D{
			{"$elemMatch", bson.D{
				{"tag", tag},
			}},
		}},
	})

	return err
}

//直接设置level等级为多少
func (self *Model) ChangeLevel(userID primitive.ObjectID, level int) error {
	_, err := self.Set(context.TODO(), self.FilterByID(userID), bson.D{
		{"level", level},
	})

	return err
}

//level正数为增加，负数为减少
func (self *Model) IncLevel(userID primitive.ObjectID, incLevel int) error {
	_, err := self.Inc(context.TODO(), self.FilterByID(userID), bson.D{
		{"level", incLevel},
	})

	return err
}

//直接设置experience经验值为多少
func (self *Model) ChangeExperience(userID primitive.ObjectID, experience int) error {
	_, err := self.Set(context.TODO(), self.FilterByID(userID), bson.D{
		{"experience", experience},
	})

	return err
}

//experience正数为增加，负数为减少
func (self *Model) IncExperience(userID primitive.ObjectID, incExperience int) error {
	_, err := self.Inc(context.TODO(), self.FilterByID(userID), bson.D{
		{"experience", incExperience},
	})

	return err
}
