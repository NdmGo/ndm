package handles

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"ndm/internal/common"
	"ndm/internal/db"
	// "ndm/internal/model"
	// "ndm/internal/op"
)

func TaskPage(c *gin.Context) {
	data := common.CommonVer()
	c.HTML(http.StatusOK, "task.tmpl", data)
}

type TaskArgs struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

func TaskList(c *gin.Context) {
	var args StoragesArgs
	if err := c.ShouldBind(&args); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}

	fmt.Println(args)

	storages, total, err := db.GetStorages(args.Page, args.Size)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessLayuiResp(c, total, "ok", storages)
}
