package model

import (
	"time"
)

type Tasks struct {
	ID       int64     `json:"id" gorm:"primaryKey"` // unique key
	MpId     int64     `json:"mp_id"`                // must be standardized
	Cron     string    `json:"cron"`                 // cron format
	Progress int64     `json:"progress"`             // progress
	LastDone string    `json:"last_done"`            // task last done time
	Modified time.Time `json:"modified"`
}
