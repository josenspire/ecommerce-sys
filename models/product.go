package models

import (
	"github.com/astaxie/beego/orm"
)

// 库存清单信息
type Inventory struct {
	InventoryId   uint64     `json:"inventoryId" orm:"column(inventoryId);PK;unique;size(64)"`
	Quantity      uint       `json:"quantity" orm:"column(quantity);default(1000)"`
	OriginPrice   float64    `json:"originPrice" orm:"column(originPrice);digits(10);decimals(2);default(1.01);"`
	SalesNum      uint       `json:"salesNum" orm:"column(salesNum);default(0)"`
	Commission    string     `json:"commission" orm:"column(commission);default(0%);description(Product commission, income=originPrice * commission)"`
	Specification string     `json:"specification" orm:"column(specification);default(默认);description(Product specification.)"`
	Status        string     `json:"status" orm:"column(status);default(inactive);on_delete(set_default);size(10)"`
	Product       []*Product `json:"product" orm:"reverse(many)"`
	BaseModel
}

type Product struct {
	ProductId      uint64     `json:"productId" orm:"column(productId);PK;unique;size(64)"`
	ProductName    string     `json:"productName" orm:"column(productName)"`
	ProductSubName string     `json:"productSubName" orm:"column(productSubName)"`
	ProductPic     string     `json:"productPic" orm:"column(productPic)"`
	ProductThum    string     `json:"productThum" orm:"column(productThum)"`
	ProductDesc    string     `json:"productDesc" orm:"column(productDesc)"`
	Priority       uint8      `json:"priority" orm:"default(0)"`
	Type           string     `json:"type" orm:"column(type);default(normal);description(Here will have 3 types, include 'recommend', 'normal', 'specific');size(20)"`
	Status         string     `json:"status" orm:"column(status);default(inactive);on_delete(set_default);size(10)"`
	Inventory      *Inventory `json:"inventory" orm:"column(inventoryId);rel(fk);"`
	BaseModel
}

type Picture struct {
	PictureId   uint64 `json:"pictureId" orm:"column(pictureId);PK;unique"`
	PictureName string `json:"pictureName" orm:"column(pictureName)"`
	Priority    uint8  `json:"priority" orm:"column(priority);default(0)"`
	Status      string `json:"status" orm:"column(status);default(inactive);on_delete(set_default);size(10)"`
	BaseModel
}

func init() {
	orm.RegisterModel(new(Inventory), new(Product), new(Picture))
}
