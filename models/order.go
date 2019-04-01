package models

import (
	"ecommerce-sys/db"
	"strings"
	"time"
)

// 订货单
type OrderForm struct {
	OrderId       uint64  `json:"orderId" gorm:"column:orderId; primary_key; not null;"`
	OrderNumber   int64   `json:"orderNumber" gorm:"column:orderNumber; not null;"`
	TotalAmount   float64 `json:"totalAmount" gorm:"column:totalAmount; type: decimal(10, 2); default: 1.00; not null;"`
	TotalDiscount float64 `json:"totalDiscount" gorm:"column:totalDiscount; type: decimal(10, 2); default: 0.00; not null;"`
	Remark        string  `json:"remark" gorm:"column:remark; not null;"`
	Invoice       string  `json:"invoice" gorm:"column:invoice; type: ENUM('PAPER', 'ELECTRONIC', 'NONE'); default:'ELECTRONIC'; not null;"`
	Status        string  `json:"status" gorm:"column:status; type: ENUM('NONPAYMENT', 'UNDERWAY', 'DELIVERED', 'COMPLETED', 'CANCEL'); default:'NONPAYMENT'; not null;"`
	OutTradeNo    uint64  `json:"outTradeNo" gorm:"column:outTradeNo; not null;"`
	UserId        uint64  `json:"userId" gorm:"column:userId; not null;"`
	AddressId     uint64  `json:"addressId" gorm:"column:addressId; not null;"`
	BaseModel
}

type Outbound struct {
	OutboundId       uint64  `json:"outboundId" gorm:"column:outboundId; primary_key; not null;"`
	ProductId        uint64  `json:"productId" gorm:"column:productId; not null;"`
	ProductName      string  `json:"productName" gorm:"column:productName; not null;"`
	ProductPic       string  `json:"productPic" gorm:"column:productPic; not null;"`
	ProductThum      string  `json:"productThum" gorm:"column:productThum; not null;"`
	ProductUnitPrice float64 `json:"productUnitPrice" gorm:"column:productUnitPrice; type: decimal(10, 2); not null;"`
	Discount         float64 `json:"discount" gorm:"column:discount; type: decimal(10, 2); not null;"`
	Count            uint8   `json:"count" gorm:"default:1; not null;"`
	Amount           float64 `json:"amount" gorm:"type: decimal(10, 2); default: 0.00; not null;"`
	Status           string  `json:"status" gorm:"column:status; default: 'active'; not null;"`
	BaseModel
}

type OrderOutbounds struct {
	ID         uint64 `json:"id" gorm:"column:id;NOT NULL;PRIMARY_KEY;"`
	OrderId    uint64 `json:"orderId" gorm:"column:orderId; not null; unique;"`
	OutboundId uint64 `json:"outboundId" gorm:"column:outboundId; not null;"`
	BaseModel
}

type OutboundDTO struct {
	ProductId        uint64  `json:"productId" gorm:"column:productId;"`
	ProductName      string  `json:"productName" gorm:"column:productName;"`
	ProductPic       string  `json:"productPic" gorm:"column:productPic;"`
	ProductThum      string  `json:"productThum" gorm:"column:productThum;"`
	ProductUnitPrice float64 `json:"productUnitPrice" gorm:"column:productUnitPrice; type: decimal(10, 2);"`
	Discount         float64 `json:"discount" gorm:"column:discount; type: decimal(10, 2);"`
	Count            uint8   `json:"count" gorm:"default:1;"`
	Amount           float64 `json:"amount" gorm:"type: decimal(10, 2); default: 0.00;"`
}

type PlaceOrderDTO struct {
	UserId        uint64  `json:"userId" gorm:"column:userId"`
	AddressId     uint64  `json:"addressId" gorm:"column:addressId"`
	TotalDiscount float64 `json:"totalDiscount" gorm:"totalDiscount; type: decimal(10, 2);"`
	TotalAmount   float64 `json:"totalAmount" gorm:"totalAmount; type: decimal(10, 2);"`
	Invoice       string  `json:"invoice" gorm:"column:invoice;"`
	Remark        string  `json:"remark" gorm:"column:remark"`
	Outbounds     []OutboundDTO
}

func (of *OrderForm) TableName() string {
	return "orderforms"
}

func (of *OrderOutbounds) TableName() string {
	return "orderoutbounds"
}

type IOrders interface {
	QueryOrders(userId uint64, orderType string, pageIndex int) (*[]OrderForm, error)
	PlaceOrder(dto *PlaceOrderDTO) error
}

func (of *OrderForm) QueryOrders(userId uint64, orderType string, pageIndex int) (*[]OrderForm, error) {
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()

	_orderType := strings.ToUpper(orderType)
	var orders []OrderForm
	if _orderType == "ALL" {
		mysqlDB.Where("userId = ?", userId).Offset((pageIndex - 1) * 20).Limit(20).Order("updatedAt DESC").Find(&orders)
	} else {
		mysqlDB.Where("userId = ? status = ?", userId, _orderType).Offset((pageIndex - 1) * 20).Limit(20).Order("updatedAt DESC").Find(&orders)
	}
	return &orders, nil
}

func (of *OrderForm) PlaceOrder(dto *PlaceOrderDTO) error {
	// mysqlDB := db.GetMySqlConnection().GetMySqlDB()

	// orderForm := buildOrderForms(dto)
	// outbounds := buildOutbounds(dto.Outbounds)
	// 	TODO
	// orderOutbounds := OrderOutbounds{
	// 	ID: GetWuid(),
	// 	OrderId: orderForm.OrderId,
	// 	OutboundId: outbounds.
	// }
}

func buildOrderForms(dto *PlaceOrderDTO) *OrderForm {
	var orderForm = OrderForm{}
	orderForm.OrderId = GetWuid()
	orderForm.OrderNumber = time.Now().Unix()
	orderForm.TotalAmount = dto.TotalAmount
	orderForm.TotalDiscount = dto.TotalDiscount
	orderForm.Remark = dto.Remark
	orderForm.Invoice = dto.Invoice
	orderForm.UserId = dto.UserId
	orderForm.AddressId = dto.AddressId
	return &orderForm
}

func buildOutbounds(dtos []OutboundDTO) *[]Outbound {
	var outbounds []Outbound
	for i, dto := range dtos {
		outbounds[i] = Outbound{
			OutboundId:       GetWuid(),
			ProductId:        dto.ProductId,
			ProductName:      dto.ProductName,
			ProductPic:       dto.ProductPic,
			ProductThum:      dto.ProductThum,
			ProductUnitPrice: dto.ProductUnitPrice,
			Discount:         dto.Discount,
			Count:            dto.Count,
			Amount:           dto.Amount,
		}
	}
	return &outbounds
}
