package op

import (
	"fmt"
	// "sort"
	"strings"
	"time"

	"ndm/internal/db"
	"ndm/internal/driver"
	"ndm/internal/errs"
	"ndm/internal/model"
	"ndm/internal/utils"
	"ndm/internal/utils/multitasking"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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

// 上传任务
func DoneTasksSync(ctx *gin.Context, mountPath string) error {
	fmt.Println("sync:", mountPath)
	storage, err := GetStorageByMountPath(mountPath)
	if err != nil {
		return err
	}

	mpid := getStoragesMpId(storage)
	fmt.Println("mpid:", mpid)
	// if multitasking.Factory(mountPath).IsRun() {
	// 	return errs.BackupTaskIsRun
	// }

	// multitasking.Factory(mountPath).SetForceRuningStatus()
	err = DoneTasksUpload(ctx, storage, mountPath)
	fmt.Println("sync:", err)
	return nil
}

func DoneTasksUpload(ctx *gin.Context, storage driver.Driver, mountPath string) error {
	// task_start := time.Now()
	root_path := getStoragesRootPath(storage)
	objs, err := StorageList(ctx, storage, "/", model.ListArgs{
		ReqPath: mountPath,
		Refresh: true,
	}, false)

	if err != nil {
		return err
	}

	for _, d := range objs {
		fpath := d.GetPath()
		fmt.Println(fpath)
		if d.IsDir() {
			relative_path := strings.ReplaceAll(fpath, root_path, "")
			fmt.Println("dir:", mountPath, fpath, relative_path)
			doneTasksUploadRecursion(ctx, storage, mountPath, relative_path)
		} else {
			if storage.GetStorage().Driver == "ftp" {
				// mtf.SetTaskLimit(1)
			}
			// mtf.DoneTask(func() {
			// 	WriteBackupLog(log_path, fpath)
			// 	err := BackupFile(ctx, storage, fpath)
			// 	if err != nil {
			// 		AddErrorLogs(err.Error())
			// 	}
			// })
		}
	}

	return nil
}

func doneTasksUploadRecursion(ctx *gin.Context, storage driver.Driver, mountPath string, path string) error {
	objs, err := StorageList(ctx, storage, path, model.ListArgs{
		ReqPath: mountPath,
		Refresh: true,
	}, false)

	// mtf := multitasking.Factory(mountPath)
	// log_path := strings.TrimPrefix(mountPath, "/")
	for _, d := range objs {
		fpath := d.GetPath()
		fmt.Println("sync d:", fpath)
		if d.IsDir() {
			doneTasksUploadRecursion(ctx, storage, mountPath, fpath)
		} else {
			if storage.GetStorage().Driver == "ftp" {
				// mtf.SetTaskLimit(1)
			}
			// mtf.DoneTask(func() {
			// 	WriteBackupLog(log_path, fpath)
			// 	err := BackupFile(ctx, storage, fpath)
			// 	if err != nil {
			// 		AddErrorLogs(err.Error())
			// 	}
			// })
		}
	}
	return err
}

// 备份任务
func DoneTasksBackup(ctx *gin.Context, mountPath string) error {
	storage, err := GetStorageByMountPath(mountPath)
	if err != nil {
		return err
	}

	if multitasking.Factory(mountPath).IsRun() {
		return errs.BackupTaskIsRun
	}

	multitasking.Factory(mountPath).SetForceRuningStatus()
	go doneTaskDownload(ctx, storage, mountPath)
	return nil
}

func doneTaskDownload(ctx *gin.Context, storage driver.Driver, mountPath string) error {
	task_start := time.Now()
	root_path := getStoragesRootPath(storage)
	objs, err := StorageList(ctx, storage, root_path, model.ListArgs{
		ReqPath: mountPath,
		Refresh: true,
	}, false)

	if err != nil {
		return err
	}

	mtf := multitasking.Factory(mountPath)
	log_path := strings.TrimPrefix(mountPath, "/")

	TruncateBackupLog(log_path)
	for _, d := range objs {
		fpath := d.GetPath()
		// fmt.Println(fpath)
		if d.IsDir() {
			doneTaskDownloadRecursion(ctx, storage, mountPath, fpath)
		} else {
			if storage.GetStorage().Driver == "ftp" {
				mtf.SetTaskLimit(1)
			}
			mtf.DoneTask(func() {
				WriteBackupLog(log_path, fpath)
				err := BackupFile(ctx, storage, fpath)
				if err != nil {
					AddErrorLogs(err.Error())
				}
			})
		}
	}

	mtf.Close()

	elapsed := time.Now().Sub(task_start)
	WriteBackupLog(log_path, fmt.Sprintf("end, cos:%s", utils.FormatDuration(elapsed)))
	return err
}

func doneTaskDownloadRecursion(ctx *gin.Context, storage driver.Driver, mountPath string, path string) error {
	objs, err := StorageList(ctx, storage, path, model.ListArgs{
		ReqPath: mountPath,
		Refresh: true,
	}, false)

	mtf := multitasking.Factory(mountPath)
	log_path := strings.TrimPrefix(mountPath, "/")
	for _, d := range objs {
		fpath := d.GetPath()

		fmt.Println(fpath)
		if d.IsDir() {
			doneTaskDownloadRecursion(ctx, storage, mountPath, fpath)
		} else {
			if storage.GetStorage().Driver == "ftp" {
				mtf.SetTaskLimit(1)
			}
			mtf.DoneTask(func() {
				WriteBackupLog(log_path, fpath)
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
