package mysql

import (
	"bluebell/models"
	"bluebell/settings"
	"testing"
)

func init() {
	// 一般使用测试数据库!!!
	cfg := &settings.MySQLConfig {
		Host: "127.0.0.1",
		User: "root",
		Password: "lzmzm1",
		Dbname: "bluebell",
		Port: 3306,
	}
	err := Init(cfg)
	if err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	p := &models.Post{
		ID: 10,
		AuthorID: 123,
		CommunityID: 1,
		Title: "test",
		Content: "just a test",
	}
	err := CreatePost(p)
	if err != nil {
		t.Fatalf("CreatPost insert record into mysql failed, err%v\n", err)
		return
	}
	t.Logf("success!")
}