package redis

import (
	"bluebell/models"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)


func GetPostIDsInOrder(p *models.ParamPostList) (ids []string, err error){
	// 1.根据order从redis中获取id
	key := getRedisKey(KeyPostTimeZset)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZset) 
	} 

	// 2.确定索引的起始点 
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1 
	// 3.从大到小查询
	ids, err = rdb.ZRevRange(ctx, key, start, end).Result()
	if err != nil {
		return nil, err
	}
	return 
}

// GetPostVoteData 查询帖子对应的赞成票数
func GetPostVoteData(ids []string) (cnt []int64, err error){
	// 法一
	// for _, id := range ids {
	// 	key := getRedisKey(KeyPostVotedZsetPrefix + id)
	// 	// 分数为1的元素的数量
	// 	val := rdb.ZCount(ctx, key, "1", "1").Val()
	// 	cnt = append(cnt, val)
	// }

	// 法二，先拼接，后执行？
	// 使用pipeline一次发送多条命令，减少RTT
	pipeline := rdb.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZsetPrefix + id)
		pipeline.ZCount(ctx, key, "1", "1")
	}

	cmders, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}

	for _, cmder := range cmders {
		val := cmder.(*redis.IntCmd).Val()
		cnt = append(cnt, val)
	}
	return
}

// GetCommunityPostIDsInOrder 根据社区ID去获取帖子ID列表
func GetCommunityPostIDsInOrder(p *models.ParamCommunityPostList)  (ids []string, err error) {
	// zinterstore 社区的帖子集合和帖子分数有序集合 

	// 确定要进行zinterstore的有序集合
	key := getRedisKey(KeyPostTimeZset)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZset) 
	} 

	// 利用缓存key减少zinterstore执行的次数
	cID := strconv.FormatInt(p.CommunityID, 10) 
	interKey := p.Order + cID 
	if rdb.Exists(ctx, interKey).Val() < 1 {
		// 不存在	
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(ctx, interKey, &redis.ZStore{
			Aggregate: "MAX",
			Keys: []string{getRedisKey(KeyCommunitySetPrefix + cID), key},
		})
		pipeline.Expire(ctx, interKey, 60*time.Second)
		_, err = pipeline.Exec(ctx)
		if err != nil {
			return nil, err 
		}
	}

	// 确定索引的起始点 
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1 
	// 3.从大到小查询
	ids, err = rdb.ZRevRange(ctx, interKey, start, end).Result()
	if err != nil {
		return nil, err
	}
	return 
}

func GetCommunityPostIDsInOrder2(p *models.ParamPostList)  (ids []string, err error) {
	// zinterstore 社区的帖子集合和帖子分数有序集合 

	// 确定要进行zinterstore的有序集合
	key := getRedisKey(KeyPostTimeZset)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZset) 
	} 

	// 利用缓存key减少zinterstore执行的次数
	cID := strconv.FormatInt(p.CommunityID, 10) 
	interKey := p.Order + cID 
	if rdb.Exists(ctx, interKey).Val() < 1 {
		// 不存在	
		pipeline := rdb.Pipeline()
		pipeline.ZInterStore(ctx, interKey, &redis.ZStore{
			Aggregate: "MAX",
			Keys: []string{getRedisKey(KeyCommunitySetPrefix + cID), key},
		})
		pipeline.Expire(ctx, interKey, 60*time.Second).Err()
		_, err = pipeline.Exec(ctx)
		if err != nil {
			return nil, err 
		}
	}

	// 确定索引的起始点 
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1 
	// 3.从大到小查询
	ids, err = rdb.ZRevRange(ctx, interKey, start, end).Result()
	if err != nil {
		return nil, err
	}
	return 
}