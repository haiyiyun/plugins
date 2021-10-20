package predefined

const (
	GroupTypeUser   = iota //群组类型：用户
	GroupTypeSystem = 999  //群组类型：系统
)

const (
	GroupStatusNormal            = iota //状态
	GroupStatusAudit                    //需审核
	GroupStatusForbidInteractive        //禁止互动
	GroupStatusForbidJoin               //禁止加入
	GroupStatusForbidQuit               //禁止退出
	GroupStatusForbidDelete             //禁止删除
)

const (
	GroupMemberStatusNormal            = iota //群组成员状态：正常
	GroupMemberStatusForbidInteractive        //禁止互动
)
