package models

import (
	"ecommerce-sys/db"
	"github.com/astaxie/beego/logs"
)

type Address struct {
	AddressId    uint64 `json:"addressId" gorm:"column:addressId; primary_key; not null;"`
	Contact      string `json:"contact" gorm:"column:contact; type:varchar(32); not null;"`
	Telephone    string `json:"telephone" gorm:"column:telephone; type:varchar(15); not null;"`
	IsDefault    bool   `json:"isDefault" gorm:"column:isDefault; default:false; not null;"`
	Country      string `json:"country" gorm:"column:country; not null;"`
	ProvinceCity string `json:"city" gorm:"column:city; not null;"`
	Details      string `json:"details" gorm:"column:details; not null;"`
	Status       string `json:"status" gorm:"column:status; type:varchar(10); default:'inactive'; not null;"`
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
	QueryAddressByAddressId(addressId string) (Address, error)
}

func (as *Address) CreateAddress(dto *AddressDTO) error {
	var address *Address
	address.AddressId = GetWuid()
	address.UserId = dto.UserId
	address.Telephone = dto.Telephone
	address.IsDefault = dto.IsDefault
	address.Country = dto.Country
	address.ProvinceCity = dto.ProvinceCity
	address.Details = dto.Details

	mysqlDB := db.GetMySqlConnection().GetMySqlDB()
	err := mysqlDB.Create(address).Error
	if err != nil {
		logs.Error(err)
	}
	return err
}

func (as *Address) QueryAddressByAddressId(addressId string) (Address, error) {
	panic("implement me")
}
