package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
)

const secret = "mingoliu"

// 将数据库操作封装成函数，等待 logic 层根据业务需求调用

// func CheckUserExist(username string) (err error) {
// 	sqlStr := `select count(user_id) from user where username = ?`
// 	var count int
// 	if err = db.Get(&count, sqlStr, username); err != nil {
// 		return  err
// 	}
// 	if count > 0 {
// 		return errors.New("用户已存在")
// 	}
// 	return 
// }

//  CheckUserExist 根据用户名检查用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err = db.Get(&count, sqlStr, username); err != nil {
		return  err		// 查询出错
	}
	if count > 0 {
		return ErrorUserExist 
	} 
	return nil
}

// func CheckUserNotExist(username string) (err error) {
// 	sqlStr := `select count(user_id) from user where username = ?`
// 	var count int
// 	if err = db.Get(&count, sqlStr, username); err != nil {
// 		return  err		// 查询出错
// 	}
// 	if count <= 0 {
// 		return ErrorUserNotExist 
// 	}
// 	return nil
// }



// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	// 对密码进行加密
	password := encryptPassword(user.Password)
	// 执行sql语句入库	
	sqlStr := `insert into user(user_id, username, password) values(?, ?, ?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, password)
	return
}

// CheckPassword 校验用户密码 
// func CheckPassword(user *models.User) (correct bool, err error) {
// 	sqlStr := `select password from user where username = ?`	
// 	var pwd string
// 	if err = db.Get(&pwd, sqlStr, user.Username); err != nil {
// 		return false, err	// 查询失败
// 	}
// 	encryptPwd := encryptPassword(user.Password)	
// 	return encryptPwd == pwd, nil 
// }


func Login(user *models.User) (err error){
	oPassword := user.Password
	sqlStr := `select user_id, username, password from user where username = ?`	
	err = db.Get(user, sqlStr, user.Username); 

	// 用户不存在
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}	
	// 判断密码是否正确
	password := encryptPassword(oPassword) 
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}

// encryptPassword 加密用户密码
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))  // 用于增强密码的安全性
	// 返回加密后的十六进制字符串
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

// GetUserByID 根据用户ID获取信息
func GetUserByID(userID int64) (user_name string, err error){		
	sqlStr := `select username from user where user_id = ?`
	if err = db.Get(&user_name, sqlStr, userID); err != nil {
		return "", err
	}
	return 
} 