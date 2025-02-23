package mysql

import (
	"bluebell/models"
	"strings"

	"github.com/jmoiron/sqlx"
)


func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post(post_id, title, content, author_id, community_id)
				values(?, ?, ?, ?, ?)`
	if _, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID); err != nil {
		return 
	}
	return
}

func GetPostDetailByID(postID int64) (data *models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
				from post where post_id = ?`	
	data = new(models.Post)
	err = db.Get(data, sqlStr, postID)				
	if err != nil {
		return 
	}
	return 
}

func GePostListByIDs(ids []string) (postList []models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
				from post 
				where post_id in (?)
				order by FIND_IN_SET(post_id, ?)
				`	
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return 
}

func GetPostList(page, size int64) (posts []models.Post, err error) {
	// 从新到旧返回帖子
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
				from post 
				order by create_time
				desc
				limit ?, ?`
	posts = make([]models.Post, 0, size)
	// (page - 1) * size 计算的是偏移量（即从哪一行开始），size 是每页要获取的行数。
	err = db.Select(&posts, sqlStr, (page - 1) * size, size)
	return
}
