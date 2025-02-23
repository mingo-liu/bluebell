package redis

import (
	"math"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// 本项目使用简化版的投票算法
// 投一票就加432分
/*
用户投票的几种情况:
direciton = 1:
	1.之前没有投过票(0)，现在投赞成票(1) ---> 更新分数和投票记录 差值绝对值：1 + 432
	2.之前投反对票(-1)	(1)								差值绝对值：2	 + 432*2
dirction = 0:
	1.之前投赞成票，现在要取消投票						 差值绝对值：1		- 432
	2.之前投反对票								    	差值绝对值：1	   +432
direciont = -1:
	1.之前没有投过票，现在投反对票				   		差值绝对值：1		- 432
	2.之前投赞成票									   差值绝对值：2       - 432*2

投票的限制：
每个帖子发表之日起的一个星期内允许用户投票
	1.到期之后，将redis中保存的票数存储到mysql表中
	2.到期之后，删除保存为文章投票的用户有序集合
*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote	= 432	// 每一票的分数
)

var (
	ErrVoteTimeExpire = errors.New("投票时间已过")
	ErrVoteRepeat	= 	errors.New("不允许重复投票")
)

func VoteForPost(userID, postID string, direct float64) (err error){
	// 1.判断投票的限制
	postTime := rdb.ZScore(ctx, getRedisKey(KeyPostTimeZset), postID).Val()
	if time.Now().Unix() - int64(postTime) > oneWeekInSeconds {
		return ErrVoteTimeExpire 
	}

	// 2和3需要放到同一个事务中
	// 2.更新帖子的分数
	// 查用户是否为帖子投过票?
	odirect := rdb.ZScore(ctx, getRedisKey(KeyPostVotedZsetPrefix + postID), userID).Val()
	if direct == odirect {	// 当前投票的类型和之前相同，提示不允许重复投票
		return ErrVoteRepeat
	}
	var op float64
	if direct > odirect {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(direct - odirect)

	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(ctx, getRedisKey(KeyPostScoreZset), diff * op * scorePerVote, postID)
	// 3.记录为帖子投过票的用户
	if direct == 0  {		// 取消投票
		pipeline.ZRem(ctx, getRedisKey(KeyPostVotedZsetPrefix + postID), userID)
	} else {
		pipeline.ZAdd(ctx, getRedisKey(KeyPostVotedZsetPrefix + postID), redis.Z{
			Score: direct,
			Member: userID,
		})
	}
	_, err = pipeline.Exec(ctx)
	return 
}

// 保存帖子的创建时间
func CreatePost(postID, cID int64) (err error){
	pipeline := rdb.TxPipeline()	// 事务操作
	// 帖子时间
	pipeline.ZAdd(ctx, getRedisKey(KeyPostTimeZset), redis.Z{
		Score: float64(time.Now().Unix()),
		Member: strconv.FormatInt(postID, 10),
	})

	// 帖子分数，以创建时间为基准，根据用户的投票增加分数
	pipeline.ZAdd(ctx, getRedisKey(KeyPostScoreZset), redis.Z{
		Score: float64(time.Now().Unix()),
		Member: strconv.FormatInt(postID, 10),
	})
	
	// 社区帖子集合
	cKey := getRedisKey(KeyCommunitySetPrefix + strconv.FormatInt(cID, 10))
	pipeline.SAdd(ctx, cKey, postID)
	_, err = pipeline.Exec(ctx)
	if err != nil {
		zap.L().Debug("redis.CreatePost failed", zap.Error(err))
		return
	}
	return
}