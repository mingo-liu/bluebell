package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	// 生成post id
	p.ID = snowflake.GenIDInt64()
	// 帖子详情入mysql	
	if err = mysql.CreatePost(p); err != nil {
		return
	}
	// 帖子创建时间入redis库
	// mysql中的和redis中的帖子创建时间不同步
	if err = redis.CreatePost(p.ID, p.CommunityID); err != nil {
		return
	}
	return
}


func GetPostDetailByID(postID int64) (postDetail *models.ApiPostDetail, err error) {
	// 查询帖子详情
	postDetail = new(models.ApiPostDetail)
	post, err := mysql.GetPostDetailByID(postID)
	if err != nil {
		zap.L().Error("mysql.GetPostDetailByID() failed",zap.Int64("post_id", postID), zap.Error(err))
		return
	}

	// 根据作者id，查询作者信息
	author_name, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID() failed", zap.Int64("auther_id", post.AuthorID), zap.Error(err))
		return
	} 

	// 根据社区id查询社区详细信息
	communityDetail, err := mysql.GetCommunityDetail(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetail() failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
		return
	}
	postDetail.Post = post
	postDetail.AuthorName = author_name 
	postDetail.CommunityDetail = communityDetail
	return 
}

func GetPostList(page, size int64) (data []models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}

	data = make([]models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		// 根据作者id，查询作者信息
		author_name, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed", zap.Int64("auther_id", post.AuthorID), zap.Error(err))
			return nil, err
		} 

		// 根据社区id查询社区详细信息
		communityDetail, err := mysql.GetCommunityDetail(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetail() failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			return nil, err
		}
		postDetail := &models.ApiPostDetail{
			AuthorName: author_name,
			Post: &post,		// post 不是指针，可以统一一下相关函数的返回值
			CommunityDetail: communityDetail,
		}
		data = append(data, *postDetail)
	}

	return
}


func GetPostList2(p *models.ParamPostList) (data []models.ApiPostDetail, err error){
	// 根据参数去redis中查询帖子ID列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return 
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder return 0 data")
		return 
	}

	// 根据ID去mysql中查询帖子的详细信息，返回的顺序和ID的顺序保持一致
	posts, err := mysql.GePostListByIDs(ids)
	if err != nil {
		return
	}

	// 提前查询好每篇帖子的投票数
	voteCount, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	data = make([]models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		// 根据作者id，查询作者信息
		author_name, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed", zap.Int64("auther_id", post.AuthorID), zap.Error(err))
			return nil, err
		} 

		// 根据社区id查询社区详细信息
		communityDetail, err := mysql.GetCommunityDetail(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetail() failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			return nil, err
		}
		postDetail := &models.ApiPostDetail{
			AuthorName: author_name,
			VoteCount: voteCount[idx],
			Post: &post,		// post 不是指针，可以统一一下相关函数的返回值
			CommunityDetail: communityDetail,
		}
		data = append(data, *postDetail)
	}
	return
}

// 融合GetPostList2 和 GetCommunityPostList2
func GetPostList3(p *models.ParamPostList) (data []models.ApiPostDetail, err error) {
	if p.CommunityID == 0 {		// 无社区ID
		data, err = GetPostList2(p)
	} else {	// 有社区ID
		data, err = GetCommunityPostList2(p)
	}
	if err != nil {
		zap.L().Error("logic:GetPostList3 failed", zap.Error(err))
		return nil, err
	}
	return
}

// GetCommunityPostList 根据社区ID返回帖子列表
func GetCommunityPostList(p *models.ParamCommunityPostList) (data []models.ApiPostDetail, err error) {
	// 根据参数去redis中查询帖子ID列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	// zap.L().Debug("redis.GetCommunityPostIDsInOrder", zap.Any("ids", ids))
	if err != nil {
		return 
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder return 0 data")
		return 
	}
	// 根据ID去mysql中查询帖子的详细信息，返回的顺序和ID的顺序保持一致
	posts, err := mysql.GePostListByIDs(ids)
	if err != nil {
		return
	}

	// 提前查询好每篇帖子的投票数
	voteCount, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	data = make([]models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		// 根据作者id，查询作者信息
		author_name, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed", zap.Int64("auther_id", post.AuthorID), zap.Error(err))
			return nil, err
		} 

		// 根据社区id查询社区详细信息
		communityDetail, err := mysql.GetCommunityDetail(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetail() failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			return nil, err
		}
		postDetail := &models.ApiPostDetail{
			AuthorName: author_name,
			VoteCount: voteCount[idx],
			Post: &post,		// post 不是指针，可以统一一下相关函数的返回值
			CommunityDetail: communityDetail,
		}
		data = append(data, *postDetail)
	}
	return
}

func GetCommunityPostList2(p *models.ParamPostList) (data []models.ApiPostDetail, err error) {
	// 根据参数去redis中查询帖子ID列表
	ids, err := redis.GetCommunityPostIDsInOrder2(p)
	if err != nil {
		return 
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder return 0 data")
		return 
	}
	// 根据ID去mysql中查询帖子的详细信息，返回的顺序和ID的顺序保持一致
	posts, err := mysql.GePostListByIDs(ids)
	if err != nil {
		return
	}

	// 提前查询好每篇帖子的投票数
	voteCount, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	data = make([]models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		// 根据作者id，查询作者信息
		author_name, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed", zap.Int64("auther_id", post.AuthorID), zap.Error(err))
			return nil, err
		} 

		// 根据社区id查询社区详细信息
		communityDetail, err := mysql.GetCommunityDetail(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetail() failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			return nil, err
		}
		postDetail := &models.ApiPostDetail{
			AuthorName: author_name,
			VoteCount: voteCount[idx],
			Post: &post,		// post 不是指针，可以统一一下相关函数的返回值
			CommunityDetail: communityDetail,
		}
		data = append(data, *postDetail)
	}
	return
}