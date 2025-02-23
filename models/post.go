package models

import "time"

// 注意内存对齐，多个tag之间加空格
type Post struct {
	ID          int64     `json:"post_id,string" db:"post_id"`
	AuthorID    int64     `json:"author_id" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"`
	Status      int64     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}


// 帖子详情
type ApiPostDetail struct {
	AuthorName  string	`json:"author_name"`	
	VoteCount	int64	`json:"vote_count"`
	*Post	
	*CommunityDetail    `json:"community"`
}
