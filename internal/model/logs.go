package model

import (
	// "fmt"
	"time"
	// "ndm/internal/utils"
)

type Logs struct {
	ID       int64     `json:"id" gorm:"primaryKey"` // unique key
	Type     string    `json:"type"`                 //
	Content  string    `json:"content"`              //
	Modified time.Time `json:"modified"`
}

func (s *Logs) GetLogs() *Logs {
	return s
}

func (s *Logs) SetLogs(logs Logs) {
	*s = logs
}
