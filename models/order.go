package models

import (
	"github.com/astaxie/beego/orm"
)

// 订货单
type OrderForm struct {
	OrderId     uint64      `json:"orderId" orm:"column(orderId);PK;unique"`
	OrderNumber uint        `json:"orderNumber" orm:"column(orderNumber)"`
	Amount      float64     `json:"amount" orm:"column(amount);digits(10);decimals(2);default(1.01)"`
	Discount    float64     `json:"discount" orm:"column(discount);digits(10);decimals(2);default(0.00)"`
	Remark      string      `json:"remark" orm:"column(remark)"`
	Invoice     string      `json:"invoice" orm:"column(invoice);default(PAGER);description(Include 'PAPER', 'ELECTRONIC', 'NONE')"`
	Status      string      `json:"status" orm:"column(status);default(inactive);on_delete(set_default);size(10)"`
	OutTradeNo  uint64      `json:"outTradeNo" orm:"column(outTradeNo);description(Order payment out trade No)"`
	UserId      uint64      `json:"userId" orm:"column(userId)"`
	AddressId   uint64      `json:"addressId" orm:"column(addressId)"`
	Outbound    []*Outbound `json:"outbound" orm:"reverse(many)"`
	BaseModel
}

type Outbound struct {
	OutboundId       uint64     `json:"outboundId" orm:"column(outboundId);PK;unique"`
	ProductId        uint64     `json:"productId" orm:"column(productId)"`
	ProductName      string     `json:"productName" orm:"column(productName)"`
	ProductPic       string     `json:"productPic" orm:"column(productPic)"`
	ProductThum      string     `json:"productThum" orm:"column(productThum)"`
	ProductUnitPrice float64    `json:"productUnitPrice" orm:"column(productUnitPrice);digits(10);decimals(2)"`
	Discount         float64    `json:"discount" orm:"column(discount);digits(10);decimals(2)"`
	Count            uint8      `json:"count" orm:"default(1)"`
	Amount           float64    `json:"amount" orm:"digits(10);decimals(2);default(0.00)"`
	Status           string     `json:"status" orm:"column(status);default(inactive);on_delete(set_default);size(10)"`
	OrderForm        *OrderForm `json:"orderForm" orm:"column(orderId);rel(fk)"`
	BaseModel
}

func (of *OrderForm) TableName() string {
	return "orderform"
}

func init() {
	orm.RegisterModel(new(OrderForm), new(Outbound))
}
