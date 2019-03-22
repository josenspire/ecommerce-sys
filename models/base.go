package models

import "time"

type BaseModel struct {
	CreatedAt time.Time  `json:"createdAt" gorm:"column:createdAt;type:dateTime"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"column:updatedAt;type:dateTime"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"column:deletedAt;type:dateTime" sql:"index"`
}

type OrmModel struct {
	ID        uint64     `json:"id" gorm:"column:id;NOT NULL;PRIMARY_KEY;"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"column:updatedAt;type:dateTime"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"column:deletedAt;type:dateTime" sql:"index"`
}
