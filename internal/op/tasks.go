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

func CreateTasks(ctx context.Context, task model.Tasks) (int64, error) {
	task.Modified = time.Now()
	task.Progress = 0
	task.Content = ""
	task.LastDone = ""
	err := db.CreateTasks(&task)
	if err != nil {
		return task.ID, errors.WithMessage(err, "failed create task in database")
	}
	return task.ID, nil
}

func DeleteTasksById(ctx context.Context, id int64) error {
	_, err := db.GetTasksById(id)
	if err != nil {
		return errors.WithMessage(err, "failed get tasks")
	}

	if err := db.DeleteTasksById(id); err != nil {
		return errors.WithMessage(err, "failed delete tasks in database")
	}
	return nil
}
