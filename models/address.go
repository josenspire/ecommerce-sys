package models

import (
	"ecommerce-sys/db"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
)

type Address struct {
	AddressId    uint64 `json:"addressId" gorm:"column:addressId; primary_key; not null;"`
	Contact      string `json:"contact" gorm:"column:contact; type:varchar(32); not null;"`
	Telephone    string `json:"telephone" gorm:"column:telephone; type:varchar(15); not null;"`
	IsDefault    bool   `json:"isDefault" gorm:"column:isDefault; default:false; not null;"`
	Country      string `json:"country" gorm:"column:country; not null;"`
	ProvinceCity string `json:"city" gorm:"column:city; not null;"`
	Details      string `json:"details" gorm:"column:details; not null;"`
	Status       string `json:"status" gorm:"column:status; type:varchar(10); default:'active'; not null;"`
	UserId       uint64 `json:"userId" gorm:"column:userId; not null;"`
	BaseModel
}

type AddressDTO struct {
	AddressId    uint64 `json:"addressId"`
	Contact      string `json:"contact"`
	Telephone    string `json:"telephone"`
	IsDefault    bool   `json:"isDefault"`
	Country      string `json:"country"`
	ProvinceCity string `json:"provinceCity"`
	Details      string `json:"details"`
	UserId       uint64 `json:"userId"`
}

type IAddress interface {
	QueryAddresses(userId uint64) (*[]Address, error)
	CreateAddress(dto *AddressDTO) error
	QueryAddressByAddressId(userId uint64, addressId uint64) (*Address, error)
	UpdateAddress(dto *AddressDTO) error
	DeleteAddressByAddressId(userId uint64, addressId uint64) error
	SetDefaultAddress(userId uint64, addressId uint64) error
}

func (addr *Address) CreateAddress(dto *AddressDTO) error {
	address := Address{}
	address.AddressId = GetWuid()
	address.Contact = dto.Contact
	address.UserId = dto.UserId
	address.Telephone = dto.Telephone
	address.IsDefault = dto.IsDefault
	address.Country = dto.Country
	address.ProvinceCity = dto.ProvinceCity
	address.Details = dto.Details

	mysqlDB := db.GetMySqlConnection().GetMySqlDB()
	ts := mysqlDB.Begin()
	if address.IsDefault == true {
		err := ts.Model(&Address{}).Where("isDefault = true and status = 'active'").Update("isDefault", false).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			logs.Error(err)
			ts.Rollback()
			return err
		}
	}
	err := ts.Create(&address).Error
	if err != nil {
		logs.Error(err)
		ts.Rollback()
	} else {
		ts.Commit()
	}
	return err
}

func (addr *Address) QueryAddressByAddressId(userId uint64, addressId uint64) (*Address, error) {
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()
	var address = Address{}

	err := mysqlDB.Where("userId = ? and addressId = ? and status = 'active'", userId, addressId).First(&address).Error
	if err != nil {
		return nil, err
	}
	return &address, nil
}

func (addr *Address) UpdateAddress(dto *AddressDTO) error {
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()
	err := mysqlDB.Model(&Address{}).Where("userId = ? and addressId = ? and status = 'active'", dto.UserId, dto.AddressId).Updates(&dto).Error
	return err
}

func (addr *Address) DeleteAddressByAddressId(userId uint64, addressId uint64) error {
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()

	mdb := mysqlDB.Where("userId = ? and addressId = ? and status = 'active'", userId, addressId).First(&Address{})
	if mdb.Error == gorm.ErrRecordNotFound {
		return mdb.Error
	}
	err := mdb.Delete(&Address{}).Error
	return err
}

func (addr *Address) SetDefaultAddress(userId uint64, addressId uint64) error {
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()

	ts := mysqlDB.Begin()
	err := ts.Model(&Address{}).Where("userId = ? and isDefault = true and status = 'active'", userId).Update("isDefault", false).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		ts.Rollback()
		logs.Error(err)
		return err
	}
	err = ts.Model(&Address{}).Where("userId = ? and addressId = ? and status = 'active'", userId, addressId).Update("isDefault", true).Error
	if err != nil {
		logs.Error(err)
		ts.Rollback()
	} else {
		ts.Commit()
	}
	return err
}

func (addr *Address) QueryAddresses(userId uint64) ([]Address, error) {
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()
	var addresses []Address

	err := mysqlDB.Where("userId = ?", userId).Order("updatedAt DESC").Find(&addresses).Error
	return addresses, err
}
