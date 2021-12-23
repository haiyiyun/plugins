package profile

import (
	"context"

	"github.com/haiyiyun/utils/help"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (self *Model) GetInfo(userID primitive.ObjectID) (pf help.M, err error) {
	filter := self.FilterByID(userID)
	filter = append(filter, bson.D{
		{"enable", true},
	}...)

	sr := self.FindOne(context.TODO(), filter)
	err = sr.Decode(&pf)

	return
}

func (self *Model) GetNickNameAndAvatar(userID primitive.ObjectID) (pf help.M, err error) {
	filter := self.FilterByID(userID)
	filter = append(filter, bson.D{
		{"enable", true},
	}...)

	opts := options.FindOne().SetProjection(bson.D{
		{"_id", 0},
		{"info.avatar", 1},
		{"info.nickname", 1},
	})

	sr := self.FindOne(context.TODO(), filter, opts)
	err = sr.Decode(&pf)

	return
}

func (self *Model) FilterByNormalProfile() bson.D {
	return bson.D{
		{"enable", true},
	}
}

func (self *Model) FilterByNickname(nickname string) bson.D {
	return bson.D{
		{"info.nickname", nickname},
	}
}

func (self *Model) FilterByNicknameWithRegex(nickname string) bson.D {
	return bson.D{
		{"info.nickname", bson.D{
			{"$regex", nickname},
			{"$options", `i`},
		}},
	}
}

func (self *Model) FilterByNicknameStartWithRegex(nickname string) bson.D {
	return bson.D{
		{"info.nickname", bson.D{
			{"$regex", `^` + nickname},
			{"$options", `im`},
		}},
	}
}
