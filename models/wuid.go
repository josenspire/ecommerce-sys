package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/edwingeng/wuid/mysql"
)

type Wuid struct {
	H uint  `gorm:"primary_key;  ; not null;"`
	X uint8 `gorm:"unique_index; default: '0'; not null;"`
}

var g *wuid.WUID

func init() {
	wuid.WithSection(10)
	g = wuid.NewWUID("default", nil)
}

func GetWuid() uint64 {
	dbUser := beego.AppConfig.String("mysqluser")
	dbPass := beego.AppConfig.String("mysqlpass")
	dbURL := beego.AppConfig.String("mysqlurls")
	dbName := beego.AppConfig.String("mysqldb")
	dbPort := beego.AppConfig.String("mysqlport")

	err := g.LoadH24FromMysql(dbURL+":"+dbPort, dbUser, dbPass, dbName, "wuids")
	if err != nil {
		logs.Error(err)
		return 0
	}
	return g.Next()
}
