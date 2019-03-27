package models

import (
	"ecommerce-sys/db"
)

type Advert struct {
	AdvertId   uint64 `json:"advertId" gorm:"column:advertId; primary_key; not null;"`
	AdvertUrl  string `json:"advertUrl" gorm:"column:advertUrl; not null;"`
	RelativeId uint64 `json:"relativeId" gorm:"column:relativeId; not null;"`
	Remark     string `json:"remark" gorm:"column:remark; not null;"`
	BaseModel
}

type IAdvert interface {
	InsertAdvert() error
	UpdateAdvertByAdvertId() error
	QueryAdvertList() ([]Advert, error)
}

func (adv *Advert) InsertAdvert() error {
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()
	adv.AdvertId = GetWuid()
	err := mysqlDB.Create(&adv).Error
	return err
}

func (adv *Advert) UpdateAdvertByAdvertId() error {
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()
	var advert = Advert{}

	err := mysqlDB.Where("advertId = ?", &adv.AdvertId).First(&advert).Error
	if err != nil {
		return err
	}
	err = mysqlDB.Model(&advert).Update(&adv).Error
	return err
}

func (adv *Advert) QueryAdvertList() ([]Advert, error) {
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()
	var advertList []Advert
	err := mysqlDB.Find(&advertList).Error
	return advertList, err
}
