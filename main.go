package main

import (
	"ecommerce-sys/models"
	_ "ecommerce-sys/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var port string

func init() {
	mysqlDBInitialize()

	// 日志：会保存手动输出的日志和系统异常日志
	// 如： logs.Error和panic
	// level 日志保存的时候的级别，默认是 Trace 级别，level值越高，记录的日志范围越广
	logs.Async()

	// perm 日志文件权限
	// filename 保存的文件名
	// maxlines 每个文件保存的最大行数，默认值 1000000
	// maxsize 每个文件保存的最大尺寸，默认值是 1 << 28, //256 MB
	// daily 是否按照每天 logrotate，默认是 true
	// maxdays 文件最多保存多少天，默认保存 7 天
	// rotate 是否开启 logrotate，默认是 true
	logs.SetLogger(logs.AdapterFile, `{"filename":"./logs/ecommerce-sys.log","level":6,"maxlines":0,"maxsize":0,"daily":true,"maxdays":30}`)
}

func mysqlDBInitialize() {
	dbUser := beego.AppConfig.String("mysqluser")
	dbPass := beego.AppConfig.String("mysqlpass")
	dbURL := beego.AppConfig.String("mysqlurls")
	dbName := beego.AppConfig.String("mysqldb")
	dbPort := beego.AppConfig.String("mysqlport")

	orm.RegisterDriver("mysql", orm.DRMySQL)
	// 参数4(可选)  设置最大空闲连接
	// 参数5(可选)  设置最大数据库连接
	maxIdle := 30
	maxConn := 30

	// 参数1        数据库的别名，用来在 ORM 中切换数据库使用
	// 参数2        driverName
	// 参数3        对应的链接字符串
	orm.RegisterDataBase(
		"default",
		"mysql",
		dbUser+":"+dbPass+"@tcp("+dbURL+":"+dbPort+")/"+dbName+"?charset=utf8&parseTime=true",
		maxIdle,
		maxConn)

	// init model
	orm.RegisterModel(new(models.User), new(models.WxSession))

	// 控制台打印查询语句
	orm.Debug = true
	// 自动建表
	orm.RunSyncdb("default", false, true)
	// 设置为 UTC 时间(default：本地时区)
	orm.DefaultTimeLoc = time.UTC
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.SetStaticPath("/image", "./static/img") // default will setup static folder, need to setup static second directory
	beego.Run()
}
