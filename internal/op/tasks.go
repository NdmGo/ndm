package op

import (
	"fmt"
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
	// "ndm/pkg/generic_sync"

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

	for _, d := range objs {
		fpath := d.GetPath()
		fmt.Println("path1:", fpath)
		if d.IsDir() {
			err := doneTaskDownloadRecursion(ctx, storage, mountPath, fpath)
			fmt.Println("DoneTaskDownloadRecursion:", err)
		} else {

			if storage.GetStorage().Driver == "ftp" {
				err := BackupFile(ctx, storage, d.GetPath())
				fmt.Println("ftp BackupFile1 err:", err)
			} else {
				multitasking.Factory(mountPath).DoneTask(func() {
					err := BackupFile(ctx, storage, fpath)
					fmt.Println("BackupFile1 err:", err)
				})
			}
		}
	}

	fmt.Println("doneTaskDownload close start")
	multitasking.Factory(mountPath).Close()
	fmt.Println("doneTaskDownload close end")
	return err
}

func doneTaskDownloadRecursion(ctx *gin.Context, storage driver.Driver, mountPath string, path string) error {
	objs, err := StorageList(ctx, storage, path, model.ListArgs{
		ReqPath: mountPath,
		Refresh: true,
	}, false)

	for _, d := range objs {
		filepath := d.GetPath()
		fmt.Println("path2:", filepath)
		if d.IsDir() {
			doneTaskDownloadRecursion(ctx, storage, mountPath, filepath)
		} else {
			if storage.GetStorage().Driver == "ftp" {
				err := BackupFile(ctx, storage, filepath)
				fmt.Println("ftp BackupFile2 err:", err)
			} else {
				fmt.Println("filepath2:", filepath)
				multitasking.Factory(mountPath).DoneTask(func() {
					err := BackupFile(ctx, storage, filepath)
					fmt.Println("BackupFile2 err:", err)
				})
			}
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
