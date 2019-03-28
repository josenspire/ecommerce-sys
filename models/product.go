package models

import (
	"ecommerce-sys/db"
	. "ecommerce-sys/utils"
	"github.com/astaxie/beego/logs"
	"strings"
)

// 库存清单信息
type Inventory struct {
	InventoryId   uint64  `json:"inventoryId" gorm:"column:inventoryId; primary_key; not null;"`
	Quantity      uint    `json:"quantity" gorm:"column:quantity; default: 1000; not null;"`
	OriginPrice   float64 `json:"originPrice" gorm:"column:originPrice; type: decimal(10, 2); default: 1.00; not null;"`
	SalesNum      uint    `json:"salesNum" gorm:"column:salesNum; default: 0; not null;"`
	Commission    string  `json:"commission" gorm:"column:commission; default: '0%'; not null;"`       // description(Product commission, income=originPrice * commission)
	Specification string  `json:"specification" gorm:"column:specification; default: '默认'; not null;"` // description(Product specification...)
	ProductId     uint64  `json:"productId" gorm:"column:productId; not null;"`
	Status        string  `json:"status" gorm:"column:status; default: 'active'; type: varchar(10); not null;"`
	BaseModel
}

type Product struct {
	ProductId      uint64     `json:"productId" gorm:"column:productId; primary_key; not null;"`
	ProductName    string     `json:"productName" gorm:"column:productName; not null;"`
	ProductSubName string     `json:"productSubName" gorm:"column:productSubName; not null;"`
	ProductPic     string     `json:"productPic" gorm:"column:productPic; not null;"`
	ProductThum    string     `json:"productThum" gorm:"column:productThum; not null;"`
	ProductDesc    string     `json:"productDesc" gorm:"column:productDesc; not null;"`
	Priority       uint8      `json:"priority" gorm:"default:0; not null;"`
	ProductType    string     `json:"productType" gorm:"column:productType; default: 'normal'; type: enum('recommend', 'normal', 'specific'); not null;"` // description(Here will have 3 types, include 'recommend', 'normal', 'specific')
	Status         string     `json:"status" gorm:"column:status; default: 'active'; type: varchar(10); not null;"`
	Inventory      *Inventory `json:"-"`
	BaseModel
}

type Picture struct {
	PictureId   uint64 `json:"pictureId" gorm:"column:pictureId; primary_key; not null;"`
	PictureName string `json:"pictureName" gorm:"column:pictureName; not null;"`
	Priority    uint8  `json:"priority" gorm:"column:priority; default: 0; not null;"`
	Status      string `json:"status" gorm:"column:status; default: 'active'; type: varchar(10); not null;"`
	BaseModel
}

type ProductDTO struct {
	Quantity      uint    `json:"quantity"`
	OriginPrice   float64 `json:"originPrice"`
	Commission    string  `json:"commission"`
	Specification string  `json:"specification"`

	ProductName    string `json:"productName"`
	ProductSubName string `json:"productSubName"`
	ProductPic     string `json:"productPic"`
	ProductThum    string `json:"productThum"`
	ProductDesc    string `json:"productDesc"`
	Priority       uint8  `json:"priority"`
	ProductType    string `json:"productType"`
}

type IProduct interface {
	InsertProduct(dtos *ProductDTO) error
	InsertMultipleProducts(dtos *[]ProductDTO) error
	QueryProductsByProductType(productType string, pageIndex int) (*[]Product, error)
	QueryProductDetails(productId uint64) (interface{}, error)
}

func (prod *Product) InsertProduct(dto *ProductDTO) error {
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()

	ts := mysqlDB.Begin()
	productModel := Product{
		ProductId:      GetWuid(),
		ProductName:    dto.ProductName,
		ProductSubName: dto.ProductSubName,
		ProductPic:     dto.ProductPic,
		ProductThum:    dto.ProductThum,
		ProductDesc:    dto.ProductDesc,
		Priority:       dto.Priority,
		ProductType:    dto.ProductType,
	}
	inventoryModel := Inventory{
		InventoryId:   GetWuid(),
		Quantity:      dto.Quantity,
		OriginPrice:   dto.OriginPrice,
		Commission:    dto.Commission,
		Specification: dto.Specification,
		ProductId:     productModel.ProductId,
	}
	err := ts.Create(&productModel).Error
	err = ts.Create(&inventoryModel).Error
	if err != nil {
		logs.Error(err)
		ts.Rollback()
	} else {
		ts.Commit()
	}
	return err
}

func (prod *Product) InsertMultipleProducts(dtos *[]ProductDTO) error {
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()

	ts := mysqlDB.Begin()
	productSqlStr, productValues := insertProducts(dtos)
	err := ts.Exec(productSqlStr, productValues...).Error

	inventorySqlStr, inventoryValues := insertInventories(dtos)
	err = ts.Exec(inventorySqlStr, inventoryValues...).Error

	if err != nil {
		logs.Error(err)
		ts.Rollback()
	} else {
		ts.Commit()
	}
	return err
}

func (prod *Product) QueryProductsByProductType(productType string, pageIndex int) (*[]Product, error) {
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()
	var products []Product

	err := mysqlDB.Where("productType = ?", productType).Offset((pageIndex - 1) * 20).Limit(20).Find(&products).Error
	return &products, err
}

func (prod *Product) QueryProductDetails(productId uint64) (interface{}, error) {
	var productDetails = make(map[string]interface{})

	var product = Product{}
	var inventories []Inventory
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()
	err := mysqlDB.Where("productId = ? and status = 'active'", productId).First(&product).Error
	err = mysqlDB.Where("productId = ? and status = 'active", productId).Find(&inventories).Error
	if err != nil {
		return nil, err
	} else {
		productDetails["productDetails"] = product
		productDetails["inventories"] = inventories
	}
	return productDetails, nil
}

func insertInventories(dtos *[]ProductDTO) (string, []interface{}) {
	var values []interface{}
	sqlStr := "INSERT INTO `inventories` (`inventoryId`, `quantity`, `originPrice`, `commission`, `specification`, `productId`, `createdAt`,`updatedAt`) VALUES"
	rowSql := "(?, ?, ?, ?, ?, ?, ?, ?)"
	var inserts []string
	for _, product := range *dtos {
		inserts = append(inserts, rowSql)
		values = append(values, GetWuid(), product.Quantity, product.OriginPrice, product.Commission, product.Specification, "productId", GenerateNowDateString(), GenerateNowDateString())
	}
	sqlStr = sqlStr + strings.Join(inserts, ",")
	return sqlStr, values
}

func insertProducts(dtos *[]ProductDTO) (string, []interface{}) {
	var values []interface{}
	sqlStr := "INSERT INTO `products` (`productId`, `productName`, `productSubName`, `productPic`, `productThum`, `productDesc`,`priority`, `productType`, `createdAt`,`updatedAt`) VALUES"
	rowSql := "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	var inserts []string
	for _, product := range *dtos {
		inserts = append(inserts, rowSql)
		values = append(values, GetWuid(), product.ProductName, product.ProductSubName, product.ProductPic, product.ProductThum, product.ProductDesc, product.Priority, product.ProductType, GenerateNowDateString(), GenerateNowDateString())
	}
	sqlStr = sqlStr + strings.Join(inserts, ",")
	return sqlStr, values
}
