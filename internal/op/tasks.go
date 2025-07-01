package op

import (
	// "fmt"
	// "sort"
	// "strings"
	"time"

	// "ndm/internal/conf"
	"ndm/internal/db"
	"ndm/internal/driver"
	"ndm/internal/errs"
	"ndm/internal/model"
	"ndm/internal/utils/multitasking"
	// "ndm/internal/utils"

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

func DoneTasksBackup(ctx *gin.Context, mountPath string) error {
	storage, err := GetStorageByMountPath(mountPath)
	if err != nil {
		return err
	}

	if multitasking.Factory(mountPath).IsRun() {
		return errs.BackupTaskIsRun
	}

	go doneTaskDownload(ctx, storage, mountPath)
	return nil
}

func doneTaskDownload(ctx *gin.Context, storage driver.Driver, mountPath string) error {
	root_path := getStoragesRootPath(storage)
	objs, err := StorageList(ctx, storage, root_path, model.ListArgs{
		ReqPath: mountPath,
		Refresh: true,
	}, false)

	if err != nil {
		return err
	}

	mtf := multitasking.Factory(mountPath)

	for _, d := range objs {
		fpath := d.GetPath()
		if d.IsDir() {
			doneTaskDownloadRecursion(ctx, storage, mountPath, fpath)
		} else {
			if storage.GetStorage().Driver == "ftp" {
				mtf.SetTaskLimit(1)
			}
			mtf.DoneTask(func() {
				err := BackupFile(ctx, storage, fpath)
				if err != nil {
					AddErrorLogs(err.Error())
				}
			})
		}
	}

	multitasking.Factory(mountPath).Close()
	return err
}

func doneTaskDownloadRecursion(ctx *gin.Context, storage driver.Driver, mountPath string, path string) error {
	objs, err := StorageList(ctx, storage, path, model.ListArgs{
		ReqPath: mountPath,
		Refresh: true,
	}, false)

	mtf := multitasking.Factory(mountPath)
	for _, d := range objs {
		fpath := d.GetPath()
		if d.IsDir() {
			doneTaskDownloadRecursion(ctx, storage, mountPath, fpath)
		} else {
			if storage.GetStorage().Driver == "ftp" {
				mtf.SetTaskLimit(1)
			}
			mtf.DoneTask(func() {
				err := BackupFile(ctx, storage, fpath)
				if err != nil {
					AddErrorLogs(err.Error())
				}
			})
		}
	}
	return err
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
