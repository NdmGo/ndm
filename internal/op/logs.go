package op

import (
	"time"

	"ndm/internal/conf"
	"ndm/internal/db"
	"ndm/internal/model"
	"ndm/internal/utils"

	"github.com/pkg/errors"
)

func AddLogs(log model.Logs) (int64, error) {
	log.Modified = time.Now()
	err := db.AddLog(&log)
	if err != nil {
		return log.ID, errors.WithMessage(err, "failed add logs in database")
	}
	return log.ID, nil
}

func AddTypeLogs(stype, content string) (int64, error) {
	var log model.Logs
	log.Type = stype
	log.Content = content
	return AddLogs(log)
}

func AddNoticeLogs(content string) (int64, error) {
	return AddTypeLogs("notice", content)
}

func AddWarnLogs(content string) (int64, error) {
	return AddTypeLogs("warn", content)
}

func AddErrorLogs(content string) (int64, error) {
	return AddTypeLogs("error", content)
}

func DeleteLogsById(id int64) error {
	// delete the logs in the database
	if err := db.DeleteLogsById(id); err != nil {
		return errors.WithMessage(err, "failed delete logs in database")
	}
	return nil
}

func TruncateLogs() error {
	// truncate the logs in the database
	if err := db.TruncateLogs(); err != nil {
		return errors.WithMessage(err, "failed truncate logs in database")
	}
	return nil
}

func TruncateBackupLog(name string) error {
	return utils.TruncateLog(conf.Log.RootPath, "backup", name)
}

func WriteBackupLog(name, content string) error {
	return utils.WriteLog(conf.Log.RootPath, "backup", name, "backup:"+content)
}

func TailBackupFile(name string, n int) ([]string, error) {
	return utils.TailFile(conf.Log.RootPath, "backup", name, n)
}

func TruncateSyncLog(name string) error {
	return utils.TruncateLog(conf.Log.RootPath, "sync", name)
}

func WriteSyncLog(name, content string) error {
	return utils.WriteLog(conf.Log.RootPath, "sync", name, "sync:"+content)
}

func TailSyncFile(name string, n int) ([]string, error) {
	return utils.TailFile(conf.Log.RootPath, "sync", name, n)
}
