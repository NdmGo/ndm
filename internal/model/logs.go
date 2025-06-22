package model

import (
	"time"
)

type Logs struct {
	ID       int64     `json:"id" gorm:"primaryKey"` // unique key
	Type     string    `json:"type"`                 //
	Content  string    `json:"content"`              //
	Modified time.Time `json:"modified"`
}
