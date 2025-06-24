package op

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"ndm/internal/db"
	"ndm/internal/model"
)

func AddLogs(ctx context.Context, log model.Logs) (int64, error) {
	log.Modified = time.Now()
	err := db.AddLog(&log)
	if err != nil {
		return log.ID, errors.WithMessage(err, "failed add logs in database")
	}
	return log.ID, nil
}

func AddTypeLogs(ctx context.Context, stype, content string) (int64, error) {
	var log model.Logs
	log.Type = stype
	log.Content = content
	return AddLogs(ctx, log)
}

func AddNoticeLogs(ctx context.Context, content string) (int64, error) {
	return AddTypeLogs(ctx, "notice", content)
}

func AddWarnLogs(ctx context.Context, content string) (int64, error) {
	return AddTypeLogs(ctx, "warn", content)
}

func AddErrorLogs(ctx context.Context, content string) (int64, error) {
	return AddTypeLogs(ctx, "error", content)
}
