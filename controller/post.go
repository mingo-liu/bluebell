package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 创建帖子
func CreatePostHandler(c *gin.Context) {
	// 1.获取参数,并校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("invalid parameters in post", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 从c中获取当前发出请求的用户id
	// 创建帖子的请求不发送用户id，Token中包含用户id
	userId, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userId

	// 2.创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3.返回响应
	ResponseSucess(c, CodeSuccess)
}

// GetPostDetailHandler 根据帖子ID获取帖子的详情
func GetPostDetailHandler(c *gin.Context) {
	// 获取帖子id
	postIDStr := c.Param("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt postID failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)	
		return
	}

	// 根据id取出帖子数据
	data, err := logic.GetPostDetailByID(postID)
	if err != nil {
		zap.L().Error("logic.GetPostDetailByID failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return	
	}

	// 返回响应
	ResponseSucess(c, data)
}


// GetPostListHandler 获取帖子列表
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	p := new(models.ParamPostList)
	page, size := GetPageInfo(c)
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("")
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 获取帖子数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 返回相应
	ResponseSucess(c, data)
}

// 升级版帖子列表接口：根据创建时间或分数返回帖子列表
func GetPostListHandler2(c *gin.Context) {
	// 1.获取参数 /api/vi/postse?page=1&size=10&order=time
	p := &models.ParamPostList{
		Page: 1,
		Size: 10,
		Order: models.OrderTime,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2.根据参数去redis中查询帖子ID列表
	// 3.根据ID去mysql中查询帖子的详细信息
	data, err := logic.GetPostList2(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 4.返回响应
	ResponseSucess(c, data)
}


// 融合GetPostListHandler2 和 GetCommunityPostListHander
func GetPostListHandler3(c *gin.Context) {
	// 1.获取参数 /api/vi/postse?page=1&size=10&order=time
	p := &models.ParamPostList{
		Page: 1,
		Size: 10,
		Order: models.OrderTime,
		CommunityID: 0,
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler3 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// zap.L().Debug("community", zap.Int64("commnuity_id", p.CommunityID))

	// 2.根据参数去redis中查询帖子ID列表
	// 3.根据ID去mysql中查询帖子的详细信息
	data, err := logic.GetPostList3(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 4.返回响应
	ResponseSucess(c, data)
}

// 根据社区类型去查找帖子列表
func GetCommunityPostListHander(c *gin.Context) {
	// 1.获取参数 /api/vi/postse?page=1&size=10&order=time
	p := &models.ParamCommunityPostList{
		ParamPostList: models.ParamPostList{
			Page: 1,
			Size: 10,
			Order: models.OrderTime,
		},
		CommunityID: 0,
	}

	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetCommunityPostListHander with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// zap.L().Debug("ParamCommunityPostList",zap.Any("cpost", p))	
	// 2.根据参数去redis中查询帖子ID列表
	// 3.根据ID去mysql中查询帖子的详细信息
	data, err := logic.GetCommunityPostList(p)
	if err != nil {
		zap.L().Error("logic.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 4.返回响应
	ResponseSucess(c, data)

}