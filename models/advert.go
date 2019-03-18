package models

import "github.com/astaxie/beego/orm"

type Advert struct {
	AdvertId   uint64 `json:"advertId" orm:"column(advertId);PK;unique;size(64)"`
	AdvertUrl  string `json:"advertUrl" orm:"column(advertUrl)"`
	RelativeId uint64 `json:"relativeId" orm:"column(relativeId)"`
	Remark     string `json:"remark" orm:"column(remark)"`
	BaseModel
}

func init() {
	orm.RegisterModel(new(Advert))
}
