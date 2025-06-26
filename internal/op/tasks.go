package op

import (
	"fmt"
	// "sort"
	// "strings"
	"time"

	// "ndm/internal/conf"
	"ndm/internal/db"
	// "ndm/internal/errs"
	"ndm/internal/driver"
	"ndm/internal/model"
	"ndm/internal/utils"
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

func printArray(arr []interface{}) {
	for i, item := range arr {
		fmt.Printf("  [%d]: ", i)

		switch v := item.(type) {
		case string:
			fmt.Printf("string: %s\n", v)
		case float64:
			fmt.Printf("number: %f\n", v)
		case map[string]interface{}:
			fmt.Println("object:")
			printMap(v)
		case []interface{}:
			fmt.Println("array:")
			printArray(v)
		default:
			fmt.Printf("unknown (%T)\n", v)
		}
	}
}

func printMap(data map[string]interface{}) {
	for key, value := range data {
		fmt.Printf("Key: %s\n", key)

		switch v := value.(type) {
		case string:
			fmt.Printf("  Type: string, Value: %s\n", v)
		case float64:
			fmt.Printf("  Type: number, Value: %f\n", v)
		case bool:
			fmt.Printf("  Type: bool, Value: %t\n", v)
		case map[string]interface{}:
			fmt.Println("  Type: object")
			printMap(v)
		case []interface{}:
			fmt.Println("  Type: array")
			printArray(v)
		case nil:
			fmt.Println("  Type: null")
		default:
			fmt.Printf("  Type: unknown (%T)\n", v)
		}
	}
}

func DoneTasksBackup(ctx *gin.Context, mountPath string) error {

	driver, err := GetStorageByMountPath(mountPath)
	root_path := getStoragesRootPath(driver)

	// if r, ok := driver.GetAddition()(driver.Addition); ok {
	// 	fmt.Println(r.BackupDir)
	// }

	addition := driver.GetAddition()

	fmt.Println("ToJson:", utils.ToJson(addition))

	// fmt.Println("ToJson:", addition["backup_dir"])

	switch v := addition.(type) {
	case map[string]interface{}:
		value := v["key"]
		fmt.Println(value)
	case map[string]string:
		value := v["key"]
		fmt.Println(value)
	default:
		fmt.Printf("不支持的类型: %T\n", v)
	}

	// 首先断言为 map[string]interface{}
	if rootMap, ok := addition.(map[string]interface{}); ok {
		printMap(rootMap)
	} else {
		fmt.Println("JSON 数据不是对象类型")
	}

	objs, err := StorageList(ctx, driver, root_path, model.ListArgs{
		ReqPath: mountPath,
		Refresh: true,
	}, false)

	for _, d := range objs {
		fmt.Println("GetPath:", d.GetPath())

		if d.IsDir() {
			DoneTaskDownloadRecursion(ctx, driver, mountPath, d.GetPath())
		} else {

		}

	}

	fmt.Println("DoneTasks[driver]:", driver, err)

	return nil
}

func DoneTaskDownloadRecursion(ctx *gin.Context, storage driver.Driver, mountPath string, path string) error {

	objs, err := StorageList(ctx, storage, path, model.ListArgs{
		ReqPath: mountPath,
		Refresh: true,
	}, false)

	for _, d := range objs {
		fmt.Println("d.GetPath():", d.GetPath())

		if d.IsDir() {
			DoneTaskDownloadRecursion(ctx, storage, mountPath, d.GetPath())
		} else {
			file, err := GetUnwrap(ctx, storage, d.GetPath())

			link, err := storage.Link(ctx, file, model.LinkArgs{
				Header:  ctx.Request.Header,
				Type:    ctx.Query("type"),
				HttpReq: ctx.Request,
			})

			storage.MakeDir(ctx, file, "/tmp")

			// link, file, err := Link(ctx, storage, d.GetPath(), model.LinkArgs{
			// 	Header:  ctx.Request.Header,
			// 	Type:    ctx.Query("type"),
			// 	HttpReq: ctx.Request,
			// })

			// fmt.Println("DoneTaskDownloadRecursion:", d.GetPath(), link, file, err)
			fmt.Println("mfile:", link.MFile, err)
			fmt.Println("url:", link.URL, err)
		}

	}

	fmt.Println("mountPath:", mountPath, "path:", path)

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
