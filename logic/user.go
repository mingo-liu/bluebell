package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

// 存放处理业务逻辑的代码

// 处理用户注册的业务逻辑
func SignUp(p *models.ParamSignup) (err error){
	// 判断用户存不存在
	// if err = mysql.CheckUserExist(p.UserName); err != nil {		// 查询出错或者用户存在
	// 	return err
	// }
	err = mysql.CheckUserExist(p.UserName) 	
	if err != nil {		// 没有错误，表明用户不存在
		return err
	}	

	// 生成UID
	userID := snowflake.GenIDInt64()
	// 构建一个User实例
	user := &models.User{
		UserID: userID,
		Username: p.UserName,
		Password: p.Password,
	}
	// 保存进数据库
	err = mysql.InsertUser(user)	
	return 
}

// 处理用户登录的业务逻辑
func Login(p *models.ParamLogin) (user *models.User, err error) {
	// 验证用户是否存在，存在则校验密码 
	user = &models.User{
		Username: p.UserName,
		Password: p.Password,
	}
	err = mysql.Login(user)
	if err != nil {
		return nil, err
	}
	// 生成jwt
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return nil, err
	}
	user.Token = token
	return 
}