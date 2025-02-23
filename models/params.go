package models

// 定义与请求参数的结构体

// 注册参数
type ParamSignup struct {
	UserName string		`json:"username" binding:"required"`	// 多个tag使用空格分开
	Password string		`json:"password" binding:"required"`
	RePassword string	`json:"re_password" binding:"required,eqfield=Password"`
}

// 登录参数
type ParamLogin struct {
	UserName string		`json:"username" binding:"required"`
	Password string		`json:"password" binding:"required"`
}


// 谁给谁投了什么票
type ParamVoteData struct {
	// UserID		// 从请求中获取用户
	PostID int64	`json:"post_id,string" binding:"required"`
	Direction int8	`json:"direction,string" binding:"oneof=-1 0 1"`	// 赞成(1) or 反对(-1) or 取消投票(0)
}

// 帖子列表参数，从URL中获取，使用form标签
type ParamPostList struct {
	Page 	int64	`json:"page" form:"page"`	
	Size 	int64	`json:"size" form:"size"`
	Order 	string	`json:"order" form:"order"`
	CommunityID 	int64 `json:"community_id" form:"community_id"`
}


// 按社区去查找帖子列表
type ParamCommunityPostList struct {
	ParamPostList
	CommunityID 	int64 `json:"community_id" form:"community_id" binding:"required"`
}

// 获取帖子列表的方式
const (
	OrderTime = "time"
	OrderScore = "score"
)