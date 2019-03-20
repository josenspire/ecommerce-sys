package db

import (
	"fmt"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"sync"
)

type MysqlConnectionPool struct{}

var mysqlInstance *MysqlConnectionPool
var once sync.Once

var db *gorm.DB
var dbErr error

func GetMySqlConnection() *MysqlConnectionPool {
	once.Do(func() {
		mysqlInstance = &MysqlConnectionPool{}
	})
	return mysqlInstance
}

func (m *MysqlConnectionPool) InitConnectionPool() bool {
	dbUser := beego.AppConfig.String("mysqluser")
	dbPass := beego.AppConfig.String("mysqlpass")
	dbURL := beego.AppConfig.String("mysqlurls")
	dbName := beego.AppConfig.String("mysqldb")
	dbPort := beego.AppConfig.String("mysqlport")

	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPass, dbURL, dbPort, dbName)
	db, dbErr = gorm.Open("mysql", uri)
	if dbErr != nil {
		log.Fatal(dbErr)
		return false
	}
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(30)

	db.LogMode(true)

	// 关闭数据库，db会被多个goroutine共享，可以不调用
	// defer db.Close()
	return true
}

func (m *MysqlConnectionPool) GetMySqlDB() (dbConnection *gorm.DB) {
	return db
}
