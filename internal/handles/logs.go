package handles

import (
	// "errors"
	// "fmt"
	"net/http"
	// "strconv"
	// "strings"

	"github.com/gin-gonic/gin"

	"ndm/internal/common"
	"ndm/internal/db"
	"ndm/internal/model"
	// "ndm/internal/op"
	// "ndm/pkg/utils"
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

	storages, total, err := db.GetStoragesDriver(args.Page, args.Size, args.Driver)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessLayuiResp(c, total, "ok", storages)
}
