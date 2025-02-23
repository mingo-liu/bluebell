package controller

import (
	"bluebell/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 获取社区列表
func CommunityHandler(c *gin.Context) {
	// 查询所有的社区 (community_id, community_name) 以列表的形式返回
	data, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)	// 不轻易将服务端报错暴露到外面
		return
	}
	ResponseSucess(c, data)		// 返回json格式的数据
}


// CommunityDetailHandler 获取社区分类详情
func CommunityDetailHandler(c *gin.Context) {
	// 1.获取社区id
	cIdStr := c.Param("id")
	cId, err := strconv.ParseInt(cIdStr, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2.获取社区id对应的详情
	data, err := logic.GetCommunityDetail(cId)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)	// 不轻易将服务端报错暴露到外面
		return
	}
	ResponseSucess(c, data)		// 返回json格式的数据
}