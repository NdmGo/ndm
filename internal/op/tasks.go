package op

import (
	"fmt"
	// "sort"
	// "io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"ndm/internal/db"
	"ndm/internal/driver"
	"ndm/internal/errs"
	"ndm/internal/model"
	"ndm/internal/stream"
	"ndm/internal/task"
	"ndm/internal/utils"
	"ndm/internal/utils/multitasking"
	pkutils "ndm/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/xhofe/tache"
)

type UpTask struct {
	task.TaskExtensionBg
	storage          driver.Driver
	dstDirActualPath string
	file             model.FileStreamer
}

func (t *UpTask) GetName() string {
	return fmt.Sprintf("upload %s to [%s](%s)", t.file.GetName(), t.storage.GetStorage().MountPath, t.dstDirActualPath)
}

func (t *UpTask) GetStatus() string {
	return "uploading"
}

func (t *UpTask) Run() error {
	t.ClearEndTime()
	t.SetStartTime(time.Now())
	defer func() { t.SetEndTime(time.Now()) }()
	return Put(t.Ctx(), t.storage, t.dstDirActualPath, t.file, t.SetProgress, true)
}

var UpTaskManager *tache.Manager[*UpTask]

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
	// fmt.Println("sync:", mountPath)
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

	dst_storage, err := GetStorageByMountPath(dst.MountPath)
	if err != nil {
		return err
	}

	// fmt.Println("mpid:", mpid)
	// if multitasking.Factory(mountPath).IsRun() {
	// 	return errs.BackupTaskIsRun
	// }

	// multitasking.Factory(mountPath).SetForceRuningStatus()
	err = DoneTasksUpload(ctx, storage, dst_storage, root_path, mountPath)
	// fmt.Println("sync:", err)
	return nil
}

func DoneTasksUpload(ctx *gin.Context, storage driver.Driver, dst_storage driver.Driver, root_path, mountPath string) error {
	// task_start := time.Now()

	objs, err := StorageList(ctx, storage, "/", model.ListArgs{
		ReqPath: mountPath,
		Refresh: true,
	}, false)

	if err != nil {
		return err
	}

	for _, d := range objs {
		fpath := d.GetPath()
		// fmt.Println(fpath)
		if d.IsDir() {
			relative_path := strings.ReplaceAll(fpath, root_path, "")
			// fmt.Println("dir:", mountPath, fpath, relative_path)
			doneTasksUploadRecursion(ctx, storage, dst_storage, root_path, mountPath, relative_path)
		} else {

			f, _ := os.Open(fpath)
			defer f.Close()

			fileInfo, err := os.Stat(fpath)
			if err != nil {
				fmt.Println("Error:", err)
			}

			fileSize := fileInfo.Size()
			modTime := fileInfo.ModTime()
			relative_path := strings.ReplaceAll(fpath, root_path, "")

			filename := filepath.Base(relative_path)
			dstDirPath := filepath.Dir(relative_path)

			mimetype := utils.GetMimeType(fpath)
			fmt.Println("file1:", mountPath, fpath, relative_path, filename, mimetype)
			// mtf.DoneTask(func() {
			// 	WriteBackupLog(log_path, fpath)
			// 	err := BackupFile(ctx, storage, fpath)
			// 	if err != nil {
			// 		AddErrorLogs(err.Error())
			// 	}
			// })

			h := make(map[*pkutils.HashType]string)

			file := &stream.FileStream{
				Obj: &model.Object{
					Name:     filename,
					Size:     fileSize,
					Modified: modTime,
					HashInfo: pkutils.NewHashInfoByMap(h),
				},
				Reader:       f,
				Mimetype:     mimetype,
				WebPutAsTask: true,
			}

			err = Put(ctx, dst_storage, dstDirPath, file, nil, true)
			fmt.Println("err1:", err)
		}
	}

	return nil
}

func doneTasksUploadRecursion(ctx *gin.Context, storage driver.Driver, dst_storage driver.Driver, root_path, mountPath string, path string) error {
	objs, err := StorageList(ctx, storage, path, model.ListArgs{
		ReqPath: mountPath,
		Refresh: true,
	}, false)

	// mtf := multitasking.Factory(mountPath)
	// log_path := strings.TrimPrefix(mountPath, "/")
	for _, d := range objs {
		fpath := d.GetPath()

		if d.IsDir() {
			doneTasksUploadRecursion(ctx, storage, dst_storage, root_path, mountPath, fpath)
		} else {

			f, _ := os.Open(fpath)
			defer f.Close()

			fileInfo, err := os.Stat(fpath)
			if err != nil {
				fmt.Println("Error:", err)
			}

			fileSize := fileInfo.Size()
			modTime := fileInfo.ModTime()
			relative_path := strings.ReplaceAll(fpath, root_path, "")

			filename := filepath.Base(relative_path)
			dstDirPath := filepath.Dir(relative_path)

			mimetype := utils.GetMimeType(fpath)
			fmt.Println("file:", mountPath, fpath, relative_path, filename, mimetype)
			// mtf.DoneTask(func() {
			// 	WriteBackupLog(log_path, fpath)
			// 	err := BackupFile(ctx, storage, fpath)
			// 	if err != nil {
			// 		AddErrorLogs(err.Error())
			// 	}
			// })

			h := make(map[*pkutils.HashType]string)

			file := &stream.FileStream{
				Obj: &model.Object{
					Name:     filename,
					Size:     fileSize,
					Modified: modTime,
					HashInfo: pkutils.NewHashInfoByMap(h),
				},
				Reader:       f,
				Mimetype:     mimetype,
				WebPutAsTask: true,
			}

			t := &UpTask{
				TaskExtensionBg: task.TaskExtensionBg{
					Creator: "bg-upload",
				},
				storage:          dst_storage,
				dstDirActualPath: relative_path,
				file:             file,
			}
			t.SetTotalBytes(file.GetSize())

			if t == nil {
				return fmt.Errorf("failed to create task")
			}

			err = Put(ctx, dst_storage, dstDirPath, file, nil, true)
			fmt.Println("err:", err)
			// fmt.Println("t:", t)
			// fmt.Println("UpTaskManager:", UpTaskManager)
			// UpTaskManager.Add(t)

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
