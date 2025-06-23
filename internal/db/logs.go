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

// GetLogsList Get all logs from database order by index
func GetLogsList(page, size int) ([]model.Logs, int64, error) {
	logDb := db.Model(&model.Logs{})
	var count int64
	if err := logDb.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get logs count")
	}

	var logs []model.Logs
	if err := db.Order(columnName("id")).Offset((page - 1) * size).Limit(size).Find(&logs).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return logs, count, nil
}
