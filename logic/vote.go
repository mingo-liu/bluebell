package logic

import (
	"bluebell/dao/redis"
	"bluebell/models"
	"strconv"

	"go.uber.org/zap"
)

// 投票功能

// 本项目使用简化版的投票算法
// 投一票就加432分

/*
用户投票的几种情况:
direciton = 1:
	1.之前没有投过票，现在投赞成票 ---> 更新分数和投票记录
	2.之前投反对票
dirction = 0:
	1.之前投赞成票，现在要取消投票
	2.之前投反对票
direciont = -1:
	1.之前没有投过票，现在投反对票
	2.之前投赞成票

投票的限制：
每个帖子发表之日起的一个星期内允许用户投票
	1.到期之后，将redis中保存的票数存储到mysql表中
	2.到期之后，删除保存为文章投票的用户有序集合
*/

// PostVote 用户为帖子投票
func VoteForPost(userID int64, p *models.ParamVoteData) (err error) {
	uID := strconv.FormatInt(userID, 10)
	pID := strconv.FormatInt(p.PostID, 10)
	zap.L().Debug("logic:VoteForPost", zap.String("userID", uID), zap.String("postID", pID), zap.Int8("direction", p.Direction))
	err = redis.VoteForPost(uID, pID, float64(p.Direction))
	return 
}