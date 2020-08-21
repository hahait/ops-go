package utils
// gorm 连接数据库初始化

import (
	"fmt"
	"time"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Db *gorm.DB
	err error
)

func dbInit() {
	dsn := "root:Opsweb1234!@(localhost)/opsweb?charset=utf8mb4&parseTime=True&loc=Local"
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		//PrepareStmt: true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		fmt.Println("数据库初始化连接失败, 错误信息: ", err)
	}

	// 开启调试模模式，可以打印出具体的 SQL 语句
	Db = Db.Debug()
	// 配置连接池
	if sql_db, err := Db.DB(); err == nil {
		sql_db.SetMaxOpenConns(100)
		sql_db.SetMaxIdleConns(10)
		sql_db.SetConnMaxLifetime(time.Hour)
	} else {
		fmt.Println("数据库连接池设置失败, 错误信息: ", err.Error())
	}
}