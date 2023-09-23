package utils

import (
	"Gous/config"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

// 连接数据库
func openDB() {
	// 获取数据库配置
	mysqlConf := config.GetGlobalConf().DbConfig
	// 连接语句
	connArgs := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		mysqlConf.User, mysqlConf.Password, mysqlConf.Host, mysqlConf.Port, mysqlConf.Dbname)
	log.Info("mdb addr: " + connArgs)

	var err error
	db, err = gorm.Open(mysql.Open(connArgs), &gorm.Config{}) // 使用默认配置连接数据库
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB() // 获取底层的 sql.DB 连接对象
	if err != nil {
		panic("fetch db connection err:" + err.Error())
	}

	sqlDB.SetMaxIdleConns(mysqlConf.MaxIdleConn)                                        // 最大空闲连接
	sqlDB.SetMaxOpenConns(mysqlConf.MaxOpenConn)                                        // 最大打开连接
	sqlDB.SetConnMaxLifetime(time.Duration(mysqlConf.MaxIdleTime * int64(time.Second))) // 最大空闲时间（s）

}

func GetDB() *gorm.DB {
	dbOnce.Do(openDB)
	return db
}
