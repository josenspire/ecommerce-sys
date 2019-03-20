package models

type Classify struct {
	ClassifyId   uint64 `json:"classifyId" gorm:"column:classifyId; not null; primary_key;"`
	ClassifyName string `json:"classifyName" gorm:"column:classifyName; type: varchar(10); not null;"`
	ClassifyIcon string `json:"classifyIcon" gorm:"column:classifyIcon; not null;"`
	Priority     uint8  `json:"priority" gorm:"column:priority; default:0; not null;"`
	Status       string `json:"status" gorm:"column:status; not null; default:'active';"`
	BaseModel
}

type Category struct {
	CateGoryId   uint64 `json:"categoryId" gorm:"column:cateGoryId; primary_key; not null;"`
	CategoryName string `json:"categoryName" gorm:"column:categoryName; not null; type:varchar(10);"`
	CategoryIcon string `json:"categoryIcon" gorm:"column:categoryIcon; not null;"`
	Priority     uint8  `json:"priority" gorm:"column:priority; not null; default:0;"`
	Status       string `json:"status" gorm:"column:status; not null; default:'active';"`
	ClassifyId   uint64 `json:"classifyId" gorm:"column:classifyId; not null;"`
	BaseModel
}
