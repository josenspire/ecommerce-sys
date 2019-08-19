package models

import (
	"database/sql"
	"ecommerce-sys/db"
	"github.com/astaxie/beego/logs"
	"github.com/edwingeng/wuid/mysql/wuid"
)

type Wuid struct {
	H uint  `gorm:"primary_key; AUTO_INCREMENT; not null;"`
	X uint8 `gorm:"unique_index; default: '0'; not null;"`
}

var g *wuid.WUID

func init() {
	wuid.WithSection(10)
	g = wuid.NewWUID("default", nil)
}

func GetWuid() uint64 {
	newDB := func() (*sql.DB, bool, error) {
		mysqlDB := db.GetMySqlConnection().GetMySqlDB()
		return mysqlDB.DB(), false, nil
	}
	// Setup
	err := g.LoadH28FromMysql(newDB, "wuids")

	if err != nil {
		logs.Error(err)
		return 0
	}
	return g.Next()
}
