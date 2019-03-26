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
	Contact      string `json:"contact"`
	Telephone    string `json:"telephone"`
	IsDefault    bool   `json:"isDefault"`
	Country      string `json:"country"`
	ProvinceCity string `json:"provinceCity"`
	Details      string `json:"details"`
	UserId       uint64 `json:"userId"`
}

type IAddress interface {
	CreateAddress(dto *AddressDTO) error
	QueryAddressByAddressId(userId uint64, addressId uint64) (Address, error)
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
	var address *Address

	err := mysqlDB.Where("userId = ? and addressId = ?", userId, addressId).First(&address).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return address, nil
}
