package db

import (
	"fmt"

	"ndm/internal/model"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func AddLog(log *model.Logs) error {
	return errors.WithStack(db.Create(log).Error)
}

// TruncateLogs just delete all logs from database
func TruncateLogs() error {
	return db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Logs{}).Error
}

// GetLogsList Get all logs from database order by index
func GetLogsList(page, size int) ([]model.Logs, int64, error) {
	logDb := db.Model(&model.Logs{})
	var count int64
	if err := logDb.Count(&count).Error; err != nil {
		return nil, 0, errors.Wrapf(err, "failed get logs count")
	}

	var logs []model.Logs

	logOrder := fmt.Sprintf("%s %s", columnName("id"), "desc")
	if err := db.Order(logOrder).Offset((page - 1) * size).Limit(size).Find(&logs).Error; err != nil {
		return nil, 0, errors.WithStack(err)
	}
	return logs, count, nil
}

// DeleteLogsById just delete logs from database by id
func DeleteLogsById(id int64) error {
	return errors.WithStack(db.Delete(&model.Logs{}, id).Error)
}
