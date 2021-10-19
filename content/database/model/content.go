package model

import (
	"time"

	"github.com/haiyiyun/mongodb/geometry"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContentGuise struct {
	Nickname string `json:"nickname" bson:"nickname" map:"nickname"`
	Avatar   string `json:"avatar" bson:"avatar" map:"avatar"`
	Sex      int    `json:"sex" bson:"sex" map:"sex"`
}

type ContentLink struct {
	External bool   `json:"external" bson:"external" map:"external"` //是否外部link
	URL      string `json:"url" bson:"url" map:"url"`                //链接地址
	Inline   bool   `json:"inline" bson:"inline" map:"inline"`       //是否内嵌，还是跳转
}

type Content struct {
	ID             primitive.ObjectID   `json:"_id" bson:"_id,omitempty" map:"_id,omitempty"`
	UserID         primitive.ObjectID   `json:"user_id" bson:"user_id" map:"user_id"`
	PublishType    int                  `json:"publish_type" bson:"publish_type" map:"publish_type"`       //发布类型，如：自己，转发等
	Type           int                  `json:"type" bson:"type" map:"type"`                               //动态类型：文章，动态等。详见：predefined.ContentTypeXXX
	AssociateType  int                  `json:"associate_type" bson:"associate_type" map:"associate_type"` //关联内容类型
	AssociateID    primitive.ObjectID   `json:"associate_id" bson:"associate_id" map:"associate_id"`       //关联内容ID
	Category       primitive.ObjectID   `json:"category" bson:"category" map:"category"`                   //分类ID
	Subject        primitive.ObjectID   `json:"subject" bson:"subject" map:"subject"`                      //主题ID
	Title          string               `json:"title" bson:"title" map:"title"`                            //标题
	AtUser         []primitive.ObjectID `json:"at_user" bson:"at_user" map:"at_user"`                      //@用户user_id
	Cover          string               `json:"cover" bson:"cover" map:"cover"`                            //封面图片
	Description    string               `json:"description" bson:"description" map:"description"`          //描述
	Author         string               `json:"author" bson:"author" map:"author"`                         //作者
	Video          string               `json:"video" bson:"video" map:"video"`                            //视频
	Voice          string               `json:"voice" bson:"voice" map:"voice"`                            //语音
	Images         []string             `json:"images" bson:"images" map:"images"`                         //图片
	Content        string               `json:"content" bson:"content" map:"content"`                      //内容，根据类型可纯文字，也可富文本
	Link           ContentLink          `json:"link" bson:"link" map:"link"`                               //链接
	Location       geometry.Point       `json:"location" bson:"location,omitempty" map:"location,omitempty"`
	Tags           []string             `json:"tags" bson:"tags" map:"tags"`                                  //标签，包括用户标签
	UserTags       []string             `json:"user_tags" bson:"user_tags" map:"user_tags"`                   //用户可编辑的标签
	Bestest        bool                 `json:"bestest" bson:"bestest" map:"bestest"`                         //是否最优
	Reliable       bool                 `json:"reliable" bson:"reliable" map:"reliable"`                      //是否靠谱
	ReadedUser     []primitive.ObjectID `json:"readed_user" bson:"readed_user" map:"readed_user"`             //阅读过用户user_id
	WantedUser     []primitive.ObjectID `json:"wanted_user" bson:"wanted_user" map:"wanted_user"`             //想要用户user_id
	LikedUser      []primitive.ObjectID `json:"liked_user" bson:"liked_user" map:"liked_user"`                //喜欢用户user_id
	HatedUser      []primitive.ObjectID `json:"hated_user" bson:"hated_user" map:"hated_user"`                //讨厌用户user_id
	Guise          ContentGuise         `json:"guise" bson:"guise" map:"guise"`                               //匿名伪装
	AntiGuiseUser  []primitive.ObjectID `json:"anti_guise_user" bson:"anti_guise_user" map:"anti_guise_user"` //反伪装的用户user_id
	Visibility     int                  `json:"visibility" bson:"visibility" map:"visibility"`                //可见度
	Copy           bool                 `json:"copy" bson:"copy" map:"copy"`                                  //是否复制，搬运。用于申明此内容非本人原创
	ForbidForward  bool                 `json:"forbid_forward" bson:"forbid_forward" map:"forbid_forward"`    //禁止转发
	ForbidDownload bool                 `json:"forbid_download" bson:"forbid_download" map:"forbid_download"` //禁止下载
	ForbidDiscuss  bool                 `json:"forbid_discuss" bson:"forbid_discuss" map:"forbid_discuss"`    //禁止评论
	Status         int                  `json:"status" bson:"status" map:"status"`
	CreateTime     time.Time            `json:"create_time" bson:"create_time" map:"create_time"`
	UpdateTime     time.Time            `json:"update_time" bson:"update_time" map:"update_time"`
}
