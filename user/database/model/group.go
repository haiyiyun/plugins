package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupMember struct {
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id" map:"user_id"`
	Adminer    bool               `bson:"adminer" json:"adminer" map:"adminer"`             //是否管理员
	Nickname   string             `bson:"nickname" json:"nickname" map:"nickname"`          //群成员昵称
	Avatar     string             `json:"avatar" bson:"avatar" map:"avatar"`                //群成员头像
	Stealth    bool               `bson:"stealth" json:"stealth" map:"stealth"`             //单独设置是否隐身，不设置隐身时，以user的设置为主
	Status     int                `bson:"status" json:"status" map:"status"`                //成员状态，如：禁止互动等，具体详见：predefined.GroupMemberStatusXXX
	JoinedTime time.Time          `bson:"joined_time" json:"joined_time" map:"joined_time"` //加入时间
}

//群组
type Group struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id" map:"_id"`
	Type          int                `bson:"type" json:"type" map:"type"`          //群组类型
	UserID        primitive.ObjectID `bson:"user_id" json:"user_id" map:"user_id"` //创建者user_id
	Name          string             `bson:"name" json:"name" map:"name"`
	Avatar        string             `json:"avatar" bson:"avatar" map:"avatar"`                   //群头像
	Introduction  string             `json:"introduction" bson:"introduction" map:"introduction"` //群介绍
	Announcement  string             `json:"announcement" bson:"announcement" map:"announcement"` //群公告
	Tags          []string           `json:"tags" bson:"tags" map:"tags"`                         //群标签
	Members       []GroupMember      `bson:"members" json:"members" map:"members"`
	Hide          bool               `bson:"hide" json:"hide" map:"hide"`                               //是否隐藏群组,隐藏后，不可查找，只能通过ID加入
	HideMembers   bool               `bson:"hide_members" json:"hide_members" map:"hide_members"`       //是否隐藏群组所有成员
	Join          bool               `bson:"join" json:"join" map:"join"`                               //是否允许加入群组
	Audit         bool               `bson:"audit" json:"audit" map:"audit"`                            //是否需要审核加入
	AuditQuestion []string           `bson:"audit_question" json:"audit_question" map:"audit_question"` //审核时需要回答的问题
	Status        int                `bson:"status" json:"status" map:"status"`                         //状态
	Delete        bool               `bson:"delete" json:"delete" map:"delete"`                         //是否删除
	DeleteTime    time.Time          `bson:"delete_time" json:"delete_time" map:"delete_time"`          //删除时间，一段删除时间后会自动删除此记录，如：3个月
	CreateTime    time.Time          `bson:"create_time" json:"create_time" map:"create_time"`          //创建时间
	UpdateTime    time.Time          `bson:"update_time" json:"update_time" map:"update_time"`          //更新时间
}
