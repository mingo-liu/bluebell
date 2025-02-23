package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)


const CtxUserIDKey = "userID"
var ErrorUserNotLogin = errors.New("用户未登录")

// GetCurrentUser 获取当前登录用户的ID
func GetCurrentUserID(c *gin.Context) (uid int64, err error) {
	userID, ok := c.Get(CtxUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return 0, err 
	}
	uid, ok = userID.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return 
	}
	return 
}


func GetPageInfo(c *gin.Context) (page, size int64) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}

	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10 
	}
	return 
}