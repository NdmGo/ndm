package handles

import (
	// "fmt"
	"net/http"
	"strconv"
	"strings"

	"ndm/internal/common"
	"ndm/internal/db"
	"ndm/internal/model"
	"ndm/internal/op"

	"github.com/gin-gonic/gin"
)

func LogsPage(c *gin.Context) {
	data := common.CommonVer()
	c.HTML(http.StatusOK, "logs.tmpl", data)
}

func LogsList(c *gin.Context) {
	var args model.PageReq
	if err := c.ShouldBind(&args); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}

	storages, total, err := db.GetLogsList(args.Page, args.Size)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessLayuiResp(c, total, "ok", storages)
}

func DeleteLogs(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if err := op.DeleteLogsById(int64(id)); err != nil {
		common.ErrorResp(c, err, 500, true)
		return
	}
	common.SuccessResp(c)
}

func TruncateLogs(c *gin.Context) {
	if err := op.TruncateLogs(); err != nil {
		common.ErrorResp(c, err, 500, true)
		return
	}
	common.SuccessResp(c)
}

type LogsReq struct {
	MountPath string `json:"mount_path" binding:"required"`
	Driver    string `json:"driver" binding:"required"`
}

func GetLogs(c *gin.Context) {
	var (
		args LogsReq
		data []string
		err  error
	)

	if err = c.ShouldBind(&args); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}

	mount_path := strings.TrimPrefix(args.MountPath, "/")

	if args.Driver == "local" {
		data, err = op.TailSyncFile(mount_path, 18)
		if err != nil {
			common.ErrorResp(c, err, 500, true)
			return
		}
	} else {
		data, err = op.TailBackupFile(mount_path, 18)
		if err != nil {
			common.ErrorResp(c, err, 500, true)
			return
		}
	}

	content := ""
	for _, d := range data {
		content += d + "\n"
	}

	common.SuccessResp(c, content)
}
