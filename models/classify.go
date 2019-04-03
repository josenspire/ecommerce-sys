package models

import (
	"ecommerce-sys/db"
	"github.com/astaxie/beego"
)

type Classify struct {
	ClassifyId       uint64 `json:"classifyId" gorm:"column:classifyId; not null; primary_key;"`
	ClassifyName     string `json:"classifyName" gorm:"column:classifyName; type: varchar(10); not null;"`
	ClassifyIcon     string `json:"classifyIcon" gorm:"column:classifyIcon; not null;"`
	ClassifyPriority uint8  `json:"classifyPriority" gorm:"column:classifyPriority; default:0; not null;"`
	Status           string `json:"status" gorm:"column:status; not null; default:'active';"`
	BaseModel
}

type Category struct {
	CategoryId       uint64 `json:"categoryId" gorm:"column:categoryId; primary_key; not null;"`
	CategoryName     string `json:"categoryName" gorm:"column:categoryName; not null; type:varchar(10);"`
	CategoryIcon     string `json:"categoryIcon" gorm:"column:categoryIcon; not null;"`
	CategoryPriority uint8  `json:"categoryPriority" gorm:"column:categoryPriority; not null; default:0;"`
	Status           string `json:"status" gorm:"column:status; not null; default:'active';"`
	ClassifyId       uint64 `json:"classifyId" gorm:"column:classifyId; not null;"`
	BaseModel
}

type ClassifyVO struct {
	ClassifyId       uint64 `json:"classifyId" gorm:"column:classifyId; not null; primary_key;"`
	ClassifyName     string `json:"classifyName" gorm:"column:classifyName; type: varchar(10); not null;"`
	ClassifyIcon     string `json:"classifyIcon" gorm:"column:classifyIcon; not null;"`
	ClassifyPriority uint8  `json:"classifyPriority" gorm:"column:classifyPriority; default:0; not null;"`

	CategoryId       uint64 `json:"categoryId" gorm:"column:categoryId; primary_key; not null;"`
	CategoryName     string `json:"categoryName" gorm:"column:categoryName; not null; type:varchar(10);"`
	CategoryIcon     string `json:"categoryIcon" gorm:"column:categoryIcon; not null;"`
	CategoryPriority uint8  `json:"categoryPriority" gorm:"column:categoryPriority; not null; default:0;"`
}

type IClassify interface {
	QueryClassifies() (*[]Classify, error)
	CreateClassify() error
	CreateCategory(category *Category) error
}

func (cls *Classify) CreateClassify() error {
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()

	cls.ClassifyId = GetWuid()
	err := mysqlDB.Create(&cls).Error
	return err
}

func (cls *Classify) CreateCategory(category *Category) error {
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()

	category.CategoryId = GetWuid()
	err := mysqlDB.Create(&category).Error
	return err
}

func (cls *Classify) QueryClassifies() ([]ClassifyVO, error) {
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()

	var classifiesVO []ClassifyVO
	err := mysqlDB.Table("classifies").Select("*").Joins("left join categories on classifies.status = 'active' and categories.status = 'active' and categories.classifyId = classifies.classifyId ORDER BY classifies.classifyPriority DESC;").Scan(&classifiesVO).Error
	if err != nil {
		beego.Error(err.Error())
		return nil, err
	}
	return classifiesVO, nil
}
