package redis

// redis key：使用命名空间的方式，方便查询和拆分

const (
	KeyPrefix              = "bluebell:"
	KeyPostTimeZset        = "post:time"   // zset: 帖子:发帖时间
	KeyPostScoreZset       = "post:score"  // zset: 帖子:帖子的分数
	KeyPostVotedZsetPrefix = "post:voted:" // zset: 用户:投票类型（赞成 or 反对）;参数是postID
	KeyCommunitySetPrefix  = "community:"  // zset: 社区:社区中的帖子
)

func getRedisKey(key string) string {
	return KeyPrefix + key
}
