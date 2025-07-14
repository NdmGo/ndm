package op

import (
	"fmt"
	// "sort"
	// "io"
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ndm/internal/db"
	"ndm/internal/driver"
	"ndm/internal/errs"
	"ndm/internal/model"
	"ndm/internal/stream"
	"ndm/internal/utils"
	"ndm/internal/utils/multitasking"
	pkutils "ndm/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func CreateTasks(t model.Tasks) (int64, error) {
	t.Modified = time.Now()
	t.Progress = 0
	t.LastDone = ""
	err := db.CreateTasks(&t)
	if err != nil {
		return t.ID, errors.WithMessage(err, "failed create task in database")
	}
	return t.ID, nil
}

// 上传任务
func DoneTasksSync(ctx *gin.Context, mountPath string) error {
	storage, err := GetStorageByMountPath(mountPath)
	if err != nil {
		return err
	}

	mpid := getStoragesMpId(storage)
	root_path := getStoragesRootPath(storage)

	dst, err := db.GetStorageById(mpid)
	if err != nil {
		return err
	}

	dstStorage, err := GetStorageByMountPath(dst.MountPath)
	if err != nil {
		return err
	}

	mtf := multitasking.Factory(mountPath, "sync")
	if mtf.IsRun() {
		return errs.BackupTaskIsRun
	}

	mtf.SetForceRuningStatus()

	if storage.GetStorage().Driver == "ftp" {
		mtf.SetTaskLimit(1)
	} else {
		mtf.SetTaskLimit(100)
	}

	go func() {
		c := context.TODO()
		err := utils.Try(c, func(c context.Context) {
			err = DoneTasksUpload(ctx, storage, dstStorage, root_path, mountPath)
		})

		if err != nil {
			log_path := strings.TrimPrefix(mountPath, "/")
			WriteBackupLog(log_path, err.Error())
		}
	}()
	return nil
}

func DoneTasksUpload(ctx *gin.Context, storage driver.Driver, dstStorage driver.Driver, root_path, mountPath string) error {
	task_start := time.Now()
	objs, err := StorageList(ctx, storage, "/", model.ListArgs{
		ReqPath: mountPath,
		Refresh: true,
	}, false)

	if err != nil {
		return err
	}

	mtf := multitasking.Factory(mountPath, "sync")
	log_path := strings.TrimPrefix(mountPath, "/")
	TruncateSyncLog(log_path)
	for _, d := range objs {
		fpath := d.GetPath()
		if d.IsDir() {
			relative_path := strings.ReplaceAll(fpath, root_path, "")
			doneTasksUploadRecursion(ctx, storage, dstStorage, root_path, mountPath, relative_path)
		} else {

			f, err := os.Open(fpath)
			if err != nil {
				fmt.Println("dd1:", err.Error())
				AddErrorLogs(err.Error())
				continue
			}
			// defer f.Close()

			fileInfo, err := os.Stat(fpath)
			if err != nil {

				fmt.Println("dd2:", err.Error())
				AddErrorLogs(err.Error())
				continue
			}

			fileSize := fileInfo.Size()
			modTime := fileInfo.ModTime()
			relative_path := strings.ReplaceAll(fpath, root_path, "")

			filename := filepath.Base(relative_path)
			dstDirPath := filepath.Dir(relative_path)

			mimetype := utils.GetMimeType(fpath)
			// fmt.Println("file1:", mountPath, fpath, relative_path, filename, mimetype)

			info := make(map[*pkutils.HashType]string)
			file := &stream.FileStream{
				Obj: &model.Object{
					Name:     filename,
					Size:     fileSize,
					Modified: modTime,
					HashInfo: pkutils.NewHashInfoByMap(info),
				},
				Reader:       f,
				Mimetype:     mimetype,
				WebPutAsTask: true,
			}

			mtf.DoneTask(func() {
				WriteSyncLog(log_path, fpath)

				err = Put(ctx, dstStorage, dstDirPath, file, nil, true)
				if err != nil {
					AddErrorLogs(err.Error())
				}
				f.Close()
			})
		}
	}

	mtf.Close()

	elapsed := time.Now().Sub(task_start)
	WriteBackupLog(log_path, fmt.Sprintf("end, cos:%s", utils.FormatDuration(elapsed)))

	return nil
}

func doneTasksUploadRecursion(ctx *gin.Context, storage driver.Driver, dstStorage driver.Driver, root_path, mountPath string, path string) error {
	objs, err := StorageList(ctx, storage, path, model.ListArgs{
		ReqPath: mountPath,
		Refresh: true,
	}, false)

	log_path := strings.TrimPrefix(mountPath, "/")
	mtf := multitasking.Factory(mountPath, "sync")

	for _, d := range objs {
		fpath := d.GetPath()

		if d.IsDir() {
			doneTasksUploadRecursion(ctx, storage, dstStorage, root_path, mountPath, fpath)
		} else {

			f, err := os.Open(fpath)
			if err != nil {
				AddErrorLogs(err.Error())
				continue
			}
			// defer f.Close()

			fileInfo, err := os.Stat(fpath)
			if err != nil {
				AddErrorLogs(err.Error())
				continue
			}

			fileSize := fileInfo.Size()
			modTime := fileInfo.ModTime()
			relative_path := strings.ReplaceAll(fpath, root_path, "")

			filename := filepath.Base(relative_path)
			dstDirPath := filepath.Dir(relative_path)

			mimetype := utils.GetMimeType(fpath)
			// fmt.Println("file:", mountPath, fpath, relative_path, filename, mimetype)
			info := make(map[*pkutils.HashType]string)
			file := &stream.FileStream{
				Obj: &model.Object{
					Name:     filename,
					Size:     fileSize,
					Modified: modTime,
					HashInfo: pkutils.NewHashInfoByMap(info),
				},
				Reader:       f,
				Mimetype:     mimetype,
				WebPutAsTask: true,
			}

			mtf.DoneTask(func() {
				WriteSyncLog(log_path, fpath)

				err = Put(ctx, dstStorage, dstDirPath, file, nil, true)
				if err != nil {
					AddErrorLogs(err.Error())
				}
				f.Close()
			})
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

	mtf := multitasking.Factory(mountPath, "backup")
	if mtf.IsRun() {
		return errs.BackupTaskIsRun
	}

	mtf.SetForceRuningStatus()

	if storage.GetStorage().Driver == "ftp" {
		mtf.SetTaskLimit(1)
	} else {
		mtf.SetTaskLimit(100)
	}

	go func() {
		c := context.TODO()
		err := utils.Try(c, func(c context.Context) {
			doneTaskDownload(ctx, storage, mountPath)
		})

		if err != nil {
			log_path := strings.TrimPrefix(mountPath, "/")
			WriteBackupLog(log_path, err.Error())
		}

	}()

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

	mtf := multitasking.Factory(mountPath, "backup")
	log_path := strings.TrimPrefix(mountPath, "/")

	TruncateBackupLog(log_path)
	for _, d := range objs {
		fpath := d.GetPath()
		if d.IsDir() {
			doneTaskDownloadRecursion(ctx, storage, mountPath, fpath)
		} else {
			// fmt.Println("doneTaskDownload file:", fpath)
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

	if err != nil {
		return err
	}

	mtf := multitasking.Factory(mountPath, "backup")
	log_path := strings.TrimPrefix(mountPath, "/")
	for _, d := range objs {
		fpath := d.GetPath()
		if d.IsDir() {
			doneTaskDownloadRecursion(ctx, storage, mountPath, fpath)
		} else {
			// fmt.Println("doneTaskDownloadRecursion file:", fpath)
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
