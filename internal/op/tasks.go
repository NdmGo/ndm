package op

import (
	"fmt"
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
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	// log "github.com/sirupsen/logrus"
)

func CreateTasks(task model.Tasks) (int64, error) {
	task.Modified = time.Now()
	task.Progress = 0
	task.LastDone = ""
	err := db.CreateTasks(&task)
	if err != nil {
		return task.ID, errors.WithMessage(err, "failed create task in database")
	}
	return task.ID, nil
}

func DoneTasks(c *gin.Context, mountPath string) error {

	driver, err := GetStorageByMountPath(mountPath)

	_objs, err := StorageList(c, driver, "/", model.ListArgs{
		ReqPath: mountPath,
		Refresh: true,
	}, false)

	for _, d := range _objs {
		fmt.Println(d.GetPath())
	}

	fmt.Println("DoneTasks:", driver, err)
	fmt.Println("DoneTasks:", _objs, err)

	return nil
}

func DeleteTasksById(id int64) error {
	_, err := db.GetTasksById(id)
	if err != nil {
		return errors.WithMessage(err, "failed get tasks")
	}

	if err := db.DeleteTasksById(id); err != nil {
		return errors.WithMessage(err, "failed delete tasks in database")
	}
	return nil
}
