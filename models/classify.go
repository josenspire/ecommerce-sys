package models

import "github.com/astaxie/beego/orm"

type Classify struct {
	ClassifyId   uint64      `json:"classifyId" orm:"column(classifyId);PK;unique;size(64)"`
	ClassifyName string      `json:"classifyName" orm:"column(classifyName);size(10)"`
	ClassifyIcon string      `json:"classifyIcon" orm:"column(classifyIcon)"`
	Priority     uint8       `json:"priority" orm:"column(priority);default(0)"`
	Status       string      `json:"status" orm:"column(status);default(inactive);on_delete(set_default)"`
	Category     []*Category `orm:"reverse(many)"`
	BaseModel
}

type Category struct {
	CateGoryId   uint64    `json:"categoryId" orm:"column(cateGoryId);PK;unique;size(64)"`
	CategoryName string    `json:"categoryName" orm:"column(categoryName);size(10)"`
	CategoryIcon string    `json:"categoryIcon" orm:"column(categoryIcon)"`
	Priority     uint8     `json:"priority" orm:"column(priority);default(0)"`
	Status       string    `json:"status" orm:"column(status);default(inactive);on_delete(set_default)"`
	Classify     *Classify `json:"classify" orm:"column(classifyId);rel(fk)"`
	BaseModel
}

func init() {
	orm.RegisterModel(new(Classify), new(Category))
}
