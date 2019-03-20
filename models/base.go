package models

import "time"

type BaseModel struct {
	CreatedAt time.Time  `gorm:"column:createdAt;type:dateTime"`
	UpdatedAt time.Time  `gorm:"column:updatedAt;type:dateTime"`
	DeletedAt *time.Time `gorm:"column:deletedAt;type:dateTime" sql:"index"`
}
