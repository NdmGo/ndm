package op

import (
	"context"
	// "fmt"
	// "sort"
	// "strings"
	"time"

	// "ndm/internal/conf"
	"ndm/internal/db"
	// "ndm/internal/errs"
	"ndm/internal/model"
	// "ndm/pkg/generic_sync"
	// "ndm/pkg/utils"
	// mapset "github.com/deckarep/golang-set/v2"
	"github.com/pkg/errors"
	// log "github.com/sirupsen/logrus"
)

func AddLogs(ctx context.Context, log model.Logs) (int64, error) {
	log.Modified = time.Now()
	err := db.AddLog(&log)
	if err != nil {
		return log.ID, errors.WithMessage(err, "failed add logs in database")
	}
	return log.ID, nil
}
