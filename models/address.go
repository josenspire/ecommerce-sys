package models

import "github.com/astaxie/beego/orm"

type Address struct {
	AddressId uint64 `json:"addressId" orm:"column(addressId);PK;unique;size(64)"`
	Contact   string `json:"contact" orm:"column(contact);size(32)"`
	Telephone string `json:"telephone" orm:"column(telephone);size(15)"`

	IsDefault bool `json:"isDefault" orm:"column(isDefault);default(false)"`

	Country      string `json:"country" orm:"column(country);null"`
	ProvinceCity string `json:"city" orm:"column(city)"`

	Status string  `json:"status" orm:"column(status);size(10);default(inactive);on_delete(set_default)"`
	User   []*User `orm:"reverse(many)"`
}

func init() {
	orm.RegisterModel(new(Address))
}
