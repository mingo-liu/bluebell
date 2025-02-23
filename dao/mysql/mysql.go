package mysql

import (
	"fmt"
	"bluebell/settings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var db *sqlx.DB

func Init(cfg *settings.MySQLConfig) (err error){
	// dsn := "root:lzmzm1@tcp(127.0.0.1:3306)/db1?charset=utf8mb4&parseTime=True"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Dbname,)
	db, err = sqlx.Connect("mysql", dsn)		// Open + Ping	
	if err != nil {
		zap.L().Error("connect DB faild", zap.Error(err))
		return
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxOpenConns)

	return
}

func Close() {
	_ = db.Close()
}
