package models

import "time"

type BaseModel struct {
	CreatedAt time.Time  `gorm:"column:createdAt"`
	UpdatedAt time.Time  `gorm:"column:updatedAt"`
	DeletedAt *time.Time `gorm:"column:deletedAt" sql:"index"`
}
