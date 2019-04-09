package main

import (
	"ecommerce-sys/db"
	"ecommerce-sys/models"
	_ "ecommerce-sys/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"
)

func init() {
	// mysqlDBInitialize()

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
	err := logs.SetLogger(logs.AdapterFile, `{"filename":"./logs/ecommerce-sys.log","level":6,"maxlines":0,"maxsize":0,"daily":true,"maxdays":30}`)
	if err != nil {
		beego.Error(err.Error())
		return
	}

	connectResult := db.GetMySqlConnection().InitConnectionPool()
	if !connectResult {
		log.Println("Init mysql database pool failure...")
		os.Exit(1)
	} else {
		log.Println("Mysql database pool init succeeded...")
	}
	initialDBTable()

	redisInit := db.GetRedisConnection().InitialRedisClient()
	if !redisInit {
		log.Println("Init redis database pool failure...")
		os.Exit(1)
	} else {
		log.Println("Redis database pool init succeeded...")
	}
}

func initialDBTable() {
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()

	if !mysqlDB.HasTable(&models.Wuid{}) {
		err := mysqlDB.Set(
			"gorm:table_options",
			"ENGINE=InnoDB DEFAULT CHARSET=utf8",
		).CreateTable(
			&models.Wuid{},
		).Error
		if err != nil {
			log.Fatal(err)
		}
	}
	if !mysqlDB.HasTable(&models.Advert{}) {
		err := mysqlDB.Set(
			"gorm:table_options",
			"ENGINE=InnoDB DEFAULT CHARSET=utf8",
		).CreateTable(
			&models.Advert{},
		).Error
		if err != nil {
			log.Fatal(err)
		}
	}
	if !mysqlDB.HasTable(&models.User{}) {
		err := mysqlDB.Set(
			"gorm:table_options",
			"ENGINE=InnoDB DEFAULT CHARSET=utf8",
		).CreateTable(
			&models.User{},
			&models.WxSession{},
			&models.Team{},
			&models.Address{},
		).Error
		if err != nil {
			log.Fatal(err)
		}
		err = mysqlDB.Model(&models.WxSession{}).AddForeignKey("userId", "users(userId)", "CASCADE", "CASCADE").Error
		err = mysqlDB.Model(&models.Team{}).AddForeignKey("userId", "users(userId)", "CASCADE", "CASCADE").Error
		err = mysqlDB.Model(&models.Address{}).AddForeignKey("userId", "users(userId)", "CASCADE", "CASCADE").Error
		if err != nil {
			log.Fatal(err)
		}
	}
	if !mysqlDB.HasTable(&models.Classify{}) {
		err := mysqlDB.Set(
			"gorm:table_options",
			"ENGINE=InnoDB DEFAULT CHARSET=utf8",
		).CreateTable(
			&models.Classify{},
			&models.Category{},
		).Error
		if err != nil {
			log.Fatal(err)
		}
		err = mysqlDB.Model(&models.Category{}).AddForeignKey("classifyId", "classifies(classifyId)", "CASCADE", "CASCADE").Error
		if err != nil {
			log.Fatal(err)
		}
	}
	if !mysqlDB.HasTable(&models.OrderForm{}) {
		err := mysqlDB.Set(
			"gorm:table_options",
			"ENGINE=InnoDB DEFAULT CHARSET=utf8",
		).CreateTable(
			&models.OrderForm{},
			&models.Outbound{},
		).Error
		if err != nil {
			log.Fatal(err)
		}
		err = mysqlDB.Model(&models.Outbound{}).AddForeignKey("orderId", "orderforms(orderId)", "CASCADE", "CASCADE").Error
		if err != nil {
			log.Fatal(err)
		}
	}
	if !mysqlDB.HasTable(&models.Inventory{}) {
		err := mysqlDB.Set(
			"gorm:table_options",
			"ENGINE=InnoDB DEFAULT CHARSET=utf8",
		).CreateTable(
			&models.Inventory{},
			&models.Product{},
			&models.Picture{},
		).Error
		if err != nil {
			log.Fatal(err)
		}
		err = mysqlDB.Model(&models.Inventory{}).AddForeignKey("productId", "Products(productId)", "CASCADE", "CASCADE").Error
		if err != nil {
			log.Fatal(err)
		}
	}
}

// beego.orm initial methods
func _() {
	dbUser := beego.AppConfig.String("mysqluser")
	dbPass := beego.AppConfig.String("mysqlpass")
	dbURL := beego.AppConfig.String("mysqlurls")
	dbName := beego.AppConfig.String("mysqldb")
	dbPort := beego.AppConfig.String("mysqlport")

	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		beego.Error(err.Error())
		return
	}
	// 参数4(可选)  设置最大空闲连接
	// 参数5(可选)  设置最大数据库连接
	maxIdle := 30
	maxConn := 30

	// 参数1        数据库的别名，用来在 ORM 中切换数据库使用
	// 参数2        driverName
	// 参数3        对应的链接字符串
	err = orm.RegisterDataBase(
		"default",
		"mysql",
		dbUser+":"+dbPass+"@tcp("+dbURL+":"+dbPort+")/"+dbName+"?charset=utf8&parseTime=true",
		maxIdle,
		maxConn)
	if err != nil {
		beego.Error(err.Error())
		return
	}

	// init model
	// orm.RegisterModel(new(WxSession), new(User), new(Address))

	// 控制台打印查询语句
	orm.Debug = true
	// 自动建表
	err = orm.RunSyncdb("default", false, true)
	if err != nil {
		beego.Error(err.Error())
		return
	}

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
