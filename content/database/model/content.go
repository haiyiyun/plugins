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
	Iframe   bool   `json:"iframe" bson:"iframe" map:"iframe"`       //是否内嵌，还是跳转到浏览器
	URL      string `json:"url" bson:"url" map:"url"`                //链接地址
}

type Content struct {
	ID                                     primitive.ObjectID   `json:"_id" bson:"_id,omitempty" map:"_id,omitempty"`
	PublishUserID                          primitive.ObjectID   `json:"publish_user_id" bson:"publish_user_id" map:"publish_user_id"`                //发布者user_id，为空代表系统发布
	Type                                   int                  `json:"type" bson:"type" map:"type"`                                                 //动态类型：文章，动态等。详见：predefined.ContentTypeXXX
	PublishType                            int                  `json:"publish_type" bson:"publish_type" map:"publish_type"`                         //发布类型，如：自己，转发等
	AssociateType                          int                  `json:"associate_type" bson:"associate_type" map:"associate_type"`                   //关联内容类型
	AssociateID                            primitive.ObjectID   `json:"associate_id" bson:"associate_id" map:"associate_id"`                         //关联内容ID
	LimitAssociateType                     int                  `json:"limit_associate_type" bson:"limit_associate_type" map:"limit_associate_type"` //限制关联类型，0为不限制（正好覆盖ContentAssociateTypeSelf），限制只能限制1种，不能多种，比如，问题，只允许回答类型来关联
	LimitAssociateNum                      int                  `json:"limit_associate_num" bson:"limit_associate_num" map:"limit_associate_num"`    //限制关联次数，0为不限制，比如，问题，只允许答应回答1次
	CategoryID                             primitive.ObjectID   `json:"category_id" bson:"category_id" map:"category_id"`                            //分类ID
	SubjectID                              primitive.ObjectID   `json:"subject_id" bson:"subject_id" map:"subject_id"`                               //主题ID
	AtUsers                                []primitive.ObjectID `json:"at_users" bson:"at_users" map:"at_users"`                                     //@用户user_id
	Author                                 string               `json:"author" bson:"author" map:"author"`                                           //作者
	Title                                  string               `json:"title" bson:"title" map:"title"`                                              //标题
	Cover                                  string               `json:"cover" bson:"cover" map:"cover"`                                              //封面图片
	Description                            string               `json:"description" bson:"description" map:"description"`                            //描述
	Video                                  string               `json:"video" bson:"video" map:"video"`                                              //视频
	Voice                                  string               `json:"voice" bson:"voice" map:"voice"`                                              //语音
	Images                                 []string             `json:"images" bson:"images" map:"images"`                                           //图片
	Content                                string               `json:"content" bson:"content" map:"content"`                                        //内容，根据类型可纯文字，也可富文本
	Link                                   ContentLink          `json:"link" bson:"link" map:"link"`                                                 //链接
	Location                               geometry.Point       `json:"location" bson:"location,omitempty" map:"location,omitempty"`
	UserTags                               []string             `json:"user_tags" bson:"user_tags" map:"user_tags"`                                                                                                             //用户可编辑的标签
	Visibility                             int                  `json:"visibility" bson:"visibility" map:"visibility"`                                                                                                          //可见度
	HideDetail                             bool                 `json:"hide_detail" bson:"hide_detail" map:"hide_detail"`                                                                                                       //隐藏详情内容
	OnlyUserIDShowDetail                   []primitive.ObjectID `json:"only_user_id_show_detail" bson:"only_user_id_show_detail" map:"only_user_id_show_detail"`                                                                //在隐藏详情内容时，只有设置的user_id可以查看详情
	Copy                                   bool                 `json:"copy" bson:"copy" map:"copy"`                                                                                                                            //是否复制，搬运。用于申明此内容非本人原创
	OnlyUserIDDiscuss                      []primitive.ObjectID `json:"only_user_id_discuss" bson:"only_user_id_discuss" map:"only_user_id_discuss"`                                                                            //只有指定的user_id可以discuss，nil为不限制，非nil影响only_publish_user_id_can_reply_discuss
	OnlyUserIDCanReplyDiscuss              []primitive.ObjectID `json:"only_publish_user_id_can_reply_discuss" bson:"only_publish_user_id_can_reply_discuss" map:"only_publish_user_id_can_reply_discuss"`                      //只有指定user_id可以回复评论，nil为不限制，only_user_id_discuss不为nil时，受only_user_id_discuss影响，需加入only_user_id_discuss
	OnlyUserIDCanNotReplyDiscuss           []primitive.ObjectID `json:"only_publish_user_id_can_not_reply_discuss" bson:"only_publish_user_id_can_not_reply_discuss" map:"only_publish_user_id_can_not_reply_discuss"`          //只有指定user_id不可以回复评论，nil为不限制，only_publish_user_id_can_not_reply_discuss权重大于only_publish_user_id_can_reply_discuss，即：only_publish_user_id_can_not_reply_discuss和only_publish_user_id_can_reply_discuss有user_id时,以only_publish_user_id_can_not_reply_discuss为准
	LimitAllDiscussNum                     int                  `json:"limit_all_discuss_num" bson:"limit_all_discuss_num" map:"limit_all_discuss_num"`                                                                         //限制所有评论次数，0为不限制，不为0时，limit_all_discuss_num控制所有的评论次数
	LimitPublishUserDiscussNum             int                  `json:"limit_publish_user_discuss_num" bson:"limit_publish_user_discuss_num" map:"limit_publish_user_discuss_num"`                                              //限制发布者评论次数，0为不限制，-1为不允许发布者评论，受limit_all_discuss_num影响
	LimitNotPublishUserAllUserDiscussNum   int                  `json:"limit_not_publish_user_all_user_discuss_num" bson:"limit_not_publish_user_all_user_discuss_num" map:"limit_not_publish_user_all_user_discuss_num"`       //限制非发布者的情况下的所有用户评论次数，0为不限制，-1为不允许非发布者评论，受limit_all_discuss_num影响
	LimitNotPublishUserEveryUserDiscussNum int                  `json:"limit_not_publish_user_every_user_discuss_num" bson:"limit_not_publish_user_every_user_discuss_num" map:"limit_not_publish_user_every_user_discuss_num"` //限制非发布者的情况下，每一个用户的评论次数，0为不限制，受limit_all_discuss_num和limit_not_publish_user_all_user_discuss_num影响
	HideDiscuss                            bool                 `json:"hide_discuss" bson:"hide_discuss" map:"hide_discuss"`                                                                                                    //隐藏评论
	OnlyUserIDShowDiscuss                  []primitive.ObjectID `json:"only_user_id_show_discuss" bson:"only_user_id_show_discuss" map:"only_user_id_show_discuss"`                                                             //在隐藏评论时，只有设置的user_id可以查看评论，只有hide_discuss为true时才有效
	ForbidForward                          bool                 `json:"forbid_forward" bson:"forbid_forward" map:"forbid_forward"`                                                                                              //禁止转发
	ForbidDownload                         bool                 `json:"forbid_download" bson:"forbid_download" map:"forbid_download"`                                                                                           //禁止下载
	ForbidDiscuss                          bool                 `json:"forbid_discuss" bson:"forbid_discuss" map:"forbid_discuss"`                                                                                              //禁止评论
	Tags                                   []string             `json:"tags" bson:"tags" map:"tags"`                                                                                                                            //标签，包括用户标签
	Bestest                                bool                 `json:"bestest" bson:"bestest" map:"bestest"`                                                                                                                   //是否最优
	Reliable                               bool                 `json:"reliable" bson:"reliable" map:"reliable"`                                                                                                                //是否靠谱
	ReadedUser                             []primitive.ObjectID `json:"readed_user" bson:"readed_user" map:"readed_user"`                                                                                                       //阅读过用户user_id
	WantedUser                             []primitive.ObjectID `json:"wanted_user" bson:"wanted_user" map:"wanted_user"`                                                                                                       //想要用户user_id
	LikedUser                              []primitive.ObjectID `json:"liked_user" bson:"liked_user" map:"liked_user"`                                                                                                          //喜欢用户user_id
	HatedUser                              []primitive.ObjectID `json:"hated_user" bson:"hated_user" map:"hated_user"`                                                                                                          //讨厌用户user_id
	Guise                                  ContentGuise         `json:"guise" bson:"guise" map:"guise"`                                                                                                                         //匿名伪装
	AntiGuiseUser                          []primitive.ObjectID `json:"anti_guise_user" bson:"anti_guise_user" map:"anti_guise_user"`                                                                                           //反伪装的用户user_id
	DiscussEstimateTotal                   int                  `json:"discuss_estimate_total" bson:"discuss_estimate_total" map:"discuss_estimate_total"`                                                                      //评论估计总数
	Status                                 int                  `json:"status" bson:"status" map:"status"`
	CreateTime                             time.Time            `json:"create_time" bson:"create_time" map:"create_time"`
	UpdateTime                             time.Time            `json:"update_time" bson:"update_time" map:"update_time"`
}
