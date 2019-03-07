package models

import (
	"github.com/astaxie/beego"
	"github.com/edwingeng/wuid/mysql"
)

type Wuid struct {
	X uint8 `json:"x" orm:"column(x);default(0)"`
	H uint  `json:"h" orm:"column(h);auto;size(10)"`
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

	g.LoadH24FromMysql(dbURL+":"+dbPort, dbUser, dbPass, dbName, "wuid")
	return g.Next()
}