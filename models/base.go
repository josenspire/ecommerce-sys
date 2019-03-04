package models

import "time"

type BaseModel struct {
	CreatedAt time.Time `json:"createdAt" orm:"column(createdAt);auto_now_add;type(datetime)"`
	UpdatedAt time.Time `json:"updatedAt" orm:"column(updatedAt);auto_now;type(datetime)"`
}
