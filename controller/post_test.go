package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	// "github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "/api/v1/post"
	r.POST(url, CreatePostHandler)

	body := `{
		"title": "test",
		"content": "just a test",
		"community_id": 2
	}`

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))

	w := httptest.NewRecorder()		// 接受路由处理器的结果
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// 方法1：判断响应内容中包不包含指定的字符串
	// assert.Contains(t, w.Body.String(), "需要登录")

	// 方法2：将响应的内容反序列化到res，然后判断每个字段与预期是否一致
	res := new(ResponseData)
	if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
		t.Fatalf("json.Unmarshal failed%v\n", err)
	}
	
	assert.Equal(t, res.Code, CodeNeedLogin)
}