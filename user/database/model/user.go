package model

import (
	"time"

	"github.com/haiyiyun/mongodb/geometry"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole struct {
	Role       string    `json:"role" bson:"role" map:"role"`
	Level      int       `json:"level" bson:"level" map:"level"`
	Experience int       `json:"experience" bson:"experience" map:"experience"`
	StartTime  time.Time `json:"start_time" bson:"start_time" map:"start_time"`
	EndTime    time.Time `json:"end_time" bson:"end_time" map:"end_time"`
}

type UserTag struct {
	Tag        string    `json:"tag" bson:"tag" map:"tag"`
	Level      int       `json:"level" bson:"level" map:"level"`
	Experience int       `json:"experience" bson:"experience" map:"experience"`
	StartTime  time.Time `json:"start_time" bson:"start_time" map:"start_time"`
	EndTime    time.Time `json:"end_time" bson:"end_time" map:"end_time"`
}

type UserOnline struct {
	Online      bool           `json:"online" bson:"online" map:"online"`
	Stealth     bool           `bson:"stealth" json:"stealth" map:"stealth"` //是否隐身
	IP          string         `json:"ip" bson:"ip" map:"ip"`
	Location    geometry.Point `json:"location" bson:"location,omitempty" map:"location,omitempty"`
	OnlineTime  time.Time      `json:"online_time" bson:"online_time" map:"online_time"`
	OfflineTime time.Time      `json:"offline_time" bson:"offline_time" map:"offline_time"`
}

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"user_id" map:"_id"`
	ExtensionID    int                `bson:"extension_id" json:"extension_id" map:"extension_id"` //扩展ID，0值代表无需和其他库同步，同步其他数据库时的自增ID，如mysql的user表的id
	Name           string             `bson:"name" json:"name" map:"name"`
	Password       string             `bson:"password" json:"-" map:"password"`
	SecurePassword string             `bson:"secure_password" json:"-" map:"secure_password"` //安全密码，提供二次验证使用
	Guest          bool               `bson:"guest" json:"guest" map:"guest"`                 //是否来宾
	Level          int                `json:"level" bson:"level" map:"level"`                 //总级别
	Experience     int                `json:"experience" bson:"experience" map:"experience"`  //总经验值
	Roles          []UserRole         `json:"roles" bson:"roles" map:"roles"`                 //用户身份角色
	Tags           []UserTag          `json:"tags" bson:"tags" map:"tags"`                    //用户身份标签
	Online         UserOnline         `json:"online" bson:"online" map:"online"`
	Location       geometry.Point     `json:"location" bson:"location,omitempty" map:"location,omitempty"`
	Enable         bool               `bson:"enable" json:"enable" map:"enable"`
	Delete         bool               `bson:"delete" json:"delete" map:"delete"`
	CreateTime     time.Time          `bson:"create_time" json:"create_time" map:"create_time"` //创建时间
	UpdateTime     time.Time          `bson:"update_time" json:"update_time" map:"update_time"` //更新时间
}
