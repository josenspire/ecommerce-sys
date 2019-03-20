package models

// 库存清单信息
type Inventory struct {
	InventoryId   uint64  `json:"inventoryId" gorm:"column:inventoryId; primary_key; not null;"`
	Quantity      uint    `json:"quantity" gorm:"column:quantity; default: 1000; not null;"`
	OriginPrice   float64 `json:"originPrice" gorm:"column:originPrice; type: decimal(10, 2); default: 1.00; not null;"`
	SalesNum      uint    `json:"salesNum" gorm:"column:salesNum; default: 0; not null;"`
	Commission    string  `json:"commission" gorm:"column:commission; default: '0%'; not null;"`       // description(Product commission, income=originPrice * commission)
	Specification string  `json:"specification" gorm:"column:specification; default: '默认'; not null;"` // description(Product specification...)
	Status        string  `json:"status" gorm:"column:status; default: 'active'; type: varchar(10); not null;"`
	BaseModel
}

type Product struct {
	ProductId      uint64 `json:"productId" gorm:"column:productId; primary_key; not null;"`
	ProductName    string `json:"productName" gorm:"column:productName; not null;"`
	ProductSubName string `json:"productSubName" gorm:"column:productSubName; not null;"`
	ProductPic     string `json:"productPic" gorm:"column:productPic; not null;"`
	ProductThum    string `json:"productThum" gorm:"column:productThum; not null;"`
	ProductDesc    string `json:"productDesc" gorm:"column:productDesc; not null;"`
	Priority       uint8  `json:"priority" gorm:"default:0; not null;"`

	Type string `json:"type" gorm:"column:type; default: 'normal'; type: enum('recommend', 'normal', 'specific'); not null;"` // description(Here will have 3 types, include 'recommend', 'normal', 'specific')

	Status      string `json:"status" gorm:"column:status; default: 'active'; type: varchar(10); not null;"`
	InventoryId uint64 `json:"inventoryId" gorm:"column:inventoryId; not null;"`
	BaseModel
}

type Picture struct {
	PictureId   uint64 `json:"pictureId" gorm:"column:pictureId; primary_key; not null;"`
	PictureName string `json:"pictureName" gorm:"column:pictureName; not null;"`
	Priority    uint8  `json:"priority" gorm:"column:priority; default: 0; not null;"`
	Status      string `json:"status" gorm:"column:status; default: 'active'; type: varchar(10); not null;"`
	BaseModel
}
