package db

import (
	// "fmt"
	// "strings"

	"github.com/pkg/errors"
	"ndm/internal/model"
)

func AddLog(log *model.Logs) error {
	return errors.WithStack(db.Create(log).Error)
}

// DeleteLogById just delete logs from database by id
func DeleteLogById(id int64) error {
	return errors.WithStack(db.Delete(&model.Logs{}, id).Error)
}

// ClearLogs just delete all logs from database
func ClearLogs(id int64) error {
	return errors.WithStack(db.Delete(&model.Logs{}, id).Error)
}
