package models


type User struct {
	UserID int64		`json:"user_id,string" db:"user_id"`
	Username string		`json:"username" db:"username"`
	Password string		`json:"password" db:"password"`
	Token 	 string 	`json:"token"`
}