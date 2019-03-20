package models

// 订货单
type OrderForm struct {
	OrderId     uint64  `json:"orderId" gorm:"column:orderId; primary_key; not null;"`
	OrderNumber uint    `json:"orderNumber" gorm:"column:orderNumber; not null;"`
	Amount      float64 `json:"amount" gorm:"column:amount; type: decimal(10, 2); default: 1.00; not null;"`
	Discount    float64 `json:"discount" gorm:"column:discount; type: decimal(10, 2); default: 0.00; not null;"`
	Remark      string  `json:"remark" gorm:"column:remark; not null;"`
	Invoice     string  `json:"invoice" gorm:"column:invoice; type: ENUM('PAPER', 'ELECTRONIC', 'NONE'); default:'PAPER'; not null;"`
	Status      string  `json:"status" gorm:"column:status; default: 'active'; type: varchar(10); not null;"`
	OutTradeNo  uint64  `json:"outTradeNo" gorm:"column:outTradeNo; not null;"`
	UserId      uint64  `json:"userId" gorm:"column:userId; not null;"`
	AddressId   uint64  `json:"addressId" gorm:"column:addressId; not null;"`
	BaseModel
}

type Outbound struct {
	OutboundId       uint64  `json:"outboundId" gorm:"column:outboundId; primary_key; not null;"`
	ProductId        uint64  `json:"productId" gorm:"column:productId not null;"`
	ProductName      string  `json:"productName" gorm:"column:productName not null;"`
	ProductPic       string  `json:"productPic" gorm:"column:productPic not null;"`
	ProductThum      string  `json:"productThum" gorm:"column:productThum; not null;"`
	ProductUnitPrice float64 `json:"productUnitPrice" gorm:"column:productUnitPrice; type: decimal(10, 2); not null;"`
	Discount         float64 `json:"discount" gorm:"column:discount; type: decimal(10, 2); not null;"`
	Count            uint8   `json:"count" gorm:"default:1; not null;"`
	Amount           float64 `json:"amount" gorm:"type: decimal(10, 2); default: 0.00; not null;"`
	Status           string  `json:"status" gorm:"column:status; default: 'active'; not null;"`
	OrderId          uint64  `json:"orderId" gorm:"column:orderId; not null;"`
	BaseModel
}

func (of *OrderForm) TableName() string {
	return "orderforms"
}
