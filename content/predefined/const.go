package predefined

const (
	PublishInterveneTypeDynamicAudit  = iota //发布干预动态类型：审核状态
	PublishInterveneTypeDynamicBlock         //发布干预动态类型：屏蔽状态
	PublishInterveneTypeDynamicDelete        //发布干预动态类型：删除状态
	PublishInterveneTypeDynamicNormal        //发布干预动态类型：正常状态

	PublishInterveneTypeArticleAudit  //发布干预文章类型：审核状态
	PublishInterveneTypeArticleBlock  //发布干预文章类型：屏蔽状态
	PublishInterveneTypeArticleDelete //发布干预文章类型：删除状态
	PublishInterveneTypeArticleNormal //发布干预文章类型：正常状态

	PublishInterveneTypeQuestionAudit  //发布干预问题类型：审核状态
	PublishInterveneTypeQuestionBlock  //发布干预问题类型：屏蔽状态
	PublishInterveneTypeQuestionDelete //发布干预问题类型：删除状态
	PublishInterveneTypeQuestionNormal //发布干预问题类型：正常状态

	PublishInterveneTypeAnswerAudit  //发布干预答案类型：审核状态
	PublishInterveneTypeAnswerBlock  //发布干预答案类型：屏蔽状态
	PublishInterveneTypeAnswerDelete //发布干预答案类型：删除状态
	PublishInterveneTypeAnswerNormal //发布干预答案类型：正常状态

	PublishInterveneTypeDiscussAudit  //发布干预评论类型：审核状态
	PublishInterveneTypeDiscussBlock  //发布干预评论类型：屏蔽状态
	PublishInterveneTypeDiscussDelete //发布干预评论类型：删除状态
	PublishInterveneTypeDiscussNormal //发布干预评论类型：正常状态

	PublishInterveneTypeMessageAudit  //发布干预消息类型：审核状态
	PublishInterveneTypeMessageBlock  //发布干预消息类型：屏蔽状态
	PublishInterveneTypeMessageDelete //发布干预消息类型：删除状态
	PublishInterveneTypeMessageNormal //发布干预消息类型：正常状态

	PublishInterveneTypeGroupMessageAudit  //发布干预群组消息类型：审核状态
	PublishInterveneTypeGroupMessageBlock  //发布干预群组消息类型：屏蔽状态
	PublishInterveneTypeGroupMessageDelete //发布干预群组消息类型：删除状态
	PublishInterveneTypeGroupMessageNormal //发布干预群组消息类型：正常状态
)

const (
	KeywordBanTypeContent = iota
	KeywordBanTypeTitle
	KeywordBanTypeSubject
	KeywordBanTypeMessage
	KeywordBanTypeGroupMessage
)

const (
	KeywordBanActionDelete  = iota //直接删除
	KeywordBanActionReplace        //直接替换
	KeywordBanActionAudit          //直接进入审核状态
	KeywordBanActionBlock          //直接进入屏蔽状态
	KeywordBanActionForbid         //直接进入禁止状态
	KeywordBanActionClose          //直接进入关闭状态
)

const (
	PublishStatusTemporary = iota //临时
	PublishStatusAudit            //审核
	PublishStatusBlock            //屏蔽
	PublishStatusForbid           //禁止所有操作
	PublishStatusFinish           //完成
	PublishStatusClose            //关闭。关闭后，不允许任何更新操作
	PublishStatusDelete           //删除
	PublishStatusNormal           //正常
)

const (
	VisibilityTypeSelf         = iota //自己
	VisibilityTypeHome                //主页
	VisibilityTypeRelationship        //关系人，如：联系人，好友，回答者等
	VisibilityTypeStranger            //陌生人
	VisibilityTypeSubject             //主题
	VisibilityTypeNearly              //附近：依赖坐标点控制
	VisibilityTypeCity                //城市：依赖坐标点控制
	VisibilityTypeProvince            //省：依赖坐标点控制
	VisibilityTypeNation              //国：依赖坐标点控制
	VisibilityTypeAll                 //所有
)

const (
	CategoryTypeSystemSelf     = iota + 1000 //系统自己类型
	CategoryTypeSystemDynamic                //系统动态类型
	CategoryTypeSystemArticle                //系统文章类型
	CategoryTypeSystemQuestion               //系统问题类型
	CategoryTypeSystemAnswer                 //系统答案类型
	CategoryTypeSharePlatform                //分享平台
)

const (
	SubjectTypeUserSelf     = iota //用户自己类型
	SubjectTypeUserDynamic         //用户动态类型
	SubjectTypeUserArticle         //用户文章类型
	SubjectTypeUserQuestion        //用户问题类型
	SubjectTypeUserAnswer          //用户答案类型
)

const (
	SubjectTypeSystemSelf     = iota + 1000 //系统自己类型
	SubjectTypeSystemDynamic                //系统动态类型
	SubjectTypeSystemArticle                //系统文章
	SubjectTypeSystemQuestion               //系统问题
	SubjectTypeSystemAnswer                 //系统答案
)

const (
	ContentPublishTypeSelf     = iota //内容发布类型：自己。类似于自己的内容或者草稿箱等
	ContentPublishTypeDynamic         //内容发布类型：动态
	ContentPublishTypeArticle         //内容发布类型：文章
	ContentPublishTypeQuestion        //内容发布类型：问题
	ContentPublishTypeAnswer          //内容发布类型：答案
)

const (
	ContentTypeVideoShort = iota //类型：短视频
	ContentTypeVideoLong         //类型：长视频
	ContentTypeVoiceShort        //类型：短语音
	ContentTypeVoiceLong         //类型：长语音
	ContentTypeImageText         //类型：图文
	ContentTypeImage             //类型：纯图片
	ContentTypeText              //类型：纯文本
)

const (
	ContentTypeSystemVideoShort = iota + 1000 //系统类型：短视频
	ContentTypeSystemVideoLong                //系统类型：长视频
	ContentTypeSystemVoiceShort               //系统类型：短语音
	ContentTypeSystemVoiceLong                //系统类型：长语音
	ContentTypeSystemImageText                //系统类型：图文
	ContentTypeSystemImage                    //系统类型：纯图片
	ContentTypeSystemText                     //系统类型：纯文本
)

const (
	ContentAssociateTypeSelf       = iota //关联类型：自己，不需要关联其他
	ContentAssociateTypeDynamic           //关联类型：动态，表示此次关联是以动态形式关联
	ContentAssociateTypeArticle           //关联类型：文章，表示此次关联是以文章形式关联
	ContentAssociateTypeQuestion          //关联类型：问题，表示此次关联是以问题形式关联
	ContentAssociateTypeAnswer            //关联类型：答案，表示此次关联是以答案形式关联
	ContentAssociateTypeCollection        //关联类型：集合，表示此次关联是以集合形式关联
	ContentAssociateTypeForward           //关联类型：转发，表示此次关联是以转发形式关联
)

const (
	ContentLinkTypeWebview = iota //普通网址
	ContentLinkTypeImage          //图片地址
	ContentLinkTypeVoice          //语音地址
	ContentLinkTypeVideo          //视频地址
)

const (
	DiscussTypeSelf     = iota //自己
	DiscussTypeDynamic         //动态
	DiscussTypeArticle         //文章
	DiscussTypeQuestion        //问题
	DiscussTypeAnswer          //答案
)

const (
	FollowTypeSelf     = iota //关注自己，自己的所有信息的变化，比如，被评论等
	FollowTypeDynamic         //关注动态，此动态下的所有信息的变化，比如，被评论等
	FollowTypeArticle         //关注文章，此文章下的所有信息的变化，比如，被评论等
	FollowTypeQuestion        //关注问题，此问题下的所有信息的变化，比如，有发表答案，问题被评论等
	FollowTypeAnswer          //关注答案，此答案下的所有信息的变化，比如，有被评论，设为最佳答案等
)

const (
	FollowTypeUser    = iota + 100 //关注用户，此用户下的所有信息的变化
	FollowTypeSubject              //关注主题，此主题下的所有信息的变化

)

const (
	FollowContentTypeSelf     = iota //关注内容类型为自己的内容
	FollowContentTypeDynamic         //关注内容类型为动态的内容
	FollowContentTypeArticle         //关注内容类型为文章的内容
	FollowContentTypeQuestion        //关注内容类型为问题的内容
	FollowContentTypeAnswer          //关注内容类型为答案的内容
)

const (
	FavoritesTypeSelf         = iota //收藏夹类型：自己
	FavoritesTypeDynamic             //收藏夹类型：动态
	FavoritesTypeArticle             //收藏夹类型：文章
	FavoritesTypeQuestion            //收藏夹类型：问题
	FavoritesTypeAnswer              //收藏夹类型：答案
	FavoritesTypeMessage             //收藏夹类型：消息
	FavoritesTypeGroupMessage        //收藏夹类型：群组消息
)

const (
	ShareTypeSelf     = iota //分享类型：自己
	ShareTypeDynamic         //分享类型：动态
	ShareTypeArticle         //分享类型：文章
	ShareTypeQuestion        //分享类型：问题
	ShareTypeAnswer          //分享类型：答案
)

const (
	ShareTypeUser = iota + 100 //分享类型：用户
)

const (
	MessageTypeUser   = iota //消息类型：用户
	MessageTypeSystem        //消息类型：系统
)

const (
	MessageStatusAudit  = iota //审核
	MessageStatusBlock         //屏蔽
	MessageStatusDelete        //删除
	MessageStatusNormal        //正常
)

const (
	MessageLinkTypeUserLink = iota //消息类型：用户链接
)

const (
	GroupMessageTypeUser   = iota //消息类型：用户
	GroupMessageTypeSystem        //消息类型：系统
)

const (
	GroupMessageStatusAudit  = iota //审核
	GroupMessageStatusBlock         //屏蔽
	GroupMessageStatusDelete        //删除
	GroupMessageStatusNormal        //正常
)

const (
	GroupMessageLinkTypeUserLink = iota //消息类型：用户链接
)
