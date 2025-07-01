package handles

import (
	// "errors"
	// "fmt"
	"net/http"
	"strconv"
	// "strings"

	"ndm/internal/common"
	"ndm/internal/db"
	"ndm/internal/model"
	"ndm/internal/op"
	// "ndm/pkg/utils"

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
