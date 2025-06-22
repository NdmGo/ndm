package model

import (
	"time"
)

type Tasks struct {
	ID   int64  `json:"id" gorm:"primaryKey"` // unique key
	Name string `json:"name"`                 //
	Type string `json:"type"`                 //
	Cron string `json:"cron"`                 //

	Modified time.Time `json:"modified"`
}
