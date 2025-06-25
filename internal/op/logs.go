package op

import (
	"time"

	"github.com/pkg/errors"
	"ndm/internal/db"
	"ndm/internal/model"
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
