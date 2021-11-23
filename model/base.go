package model

import (
	"apollo-proxy/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

// 删除状态
const (
	NORMOL = iota
	DELETE
)

// 任务状态
const (
	TASK_PENDING = iota + 1
	TASK_SUCCESS
	TASK_ABORT
	TASK_FAIL
	TASK_NOT_CREATE
)
const MAX_RETRY = 20

type BaseModel struct {
}

var DB *gorm.DB

func init() {
	DB = GetDB()
	if DB == nil {
		log.Fatalln("数据库连接创建失败")
	}
}

func GetDB() *gorm.DB {
	dsn := config.Config.Mysql.User + ":" + config.Config.Mysql.Passwd +
		"@tcp(" + config.Config.Mysql.Host + ")/" + config.Config.Mysql.DbName +
		"?multiStatements=true&parseTime=true&charset=utf8mb4&loc=Asia%2FShanghai"
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Println("数据库连接创建失败: ", err.Error(), dsn)
		return nil
	}
	//db.LogMode(common.GetBool("application.debug", false))
	db.LogMode(true)

	return db
}
