package models

import "time"

type BaseModel struct {
	CreatedAt time.Time  `json:"createdAt" gorm:"column:createdAt;type:dateTime"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"column:updatedAt;type:dateTime"`
	DeletedAt *time.Time `json:"-" gorm:"column:deletedAt;type:dateTime" sql:"index"`
}
