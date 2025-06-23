package model

import (
	"time"
)

type Tasks struct {
	ID        int64     `json:"id" gorm:"primaryKey"`     // unique key
	MountPath string    `json:"mount_path" gorm:"unique"` // must be standardized
	Cron      string    `json:"cron"`                     //
	Progress  int64     `json:"progress"`                 //
	Content   string    `json:"content"`                  //
	Modified  time.Time `json:"modified"`
}
