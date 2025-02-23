package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// 处理用户的注册请求
func SignUpHandler(c *gin.Context) {
	// 1.获取参数和参数校验
	p := new(models.ParamSignup)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 验证是否是validator中的校验错误
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)	
			return 
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))	
		return
	}
	// fmt.Println(*p)

	// 2.业务处理
	if err := logic.SignUp(p); err != nil {
		if errors.Is(err, mysql.ErrorUserExist){
			ResponseError(c, CodeUserExist)
			return 
		}
		ResponseError(c, CodeServerBusy)
		return 
	}

	// 3.返回响应
	ResponseSucess(c, nil)
}

// LoginHandler 处理用户登录请求
func LoginHandler(c *gin.Context) {
	// 获取参数并校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		// 验证是否是validator中的校验错误
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)	
			return 
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))	
		return
	}

	// 业务处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed....", zap.String("username", p.UserName), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)	// 用户不存在
			return
		}
		ResponseError(c, CodeInvalidPassword)	// 用户名或密码错误
		// 还有一种查询失败的err
		return 
	}
	// 返回响应
	ResponseSucess(c, user)	
}