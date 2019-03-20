package models

type Advert struct {
	AdvertId   uint64 `json:"advertId" gorm:"column:advertId; primary_key; not null;"`
	AdvertUrl  string `json:"advertUrl" gorm:"column:advertUrl; not null;"`
	RelativeId uint64 `json:"relativeId" gorm:"column:relativeId; not null;"`
	Remark     string `json:"remark" gorm:"column:remark; not null;"`
	BaseModel
}
