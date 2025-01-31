package handles

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"ndm/internal/common"
	"ndm/internal/db"
	"ndm/internal/model"
	// "ndm/internal/op"
)

func StoragesPage(c *gin.Context) {
	data := common.CommonVer()
	c.HTML(http.StatusOK, "storage.tmpl", data)
}

func StoragesEditPage(c *gin.Context) {
	data := common.CommonVer()
	c.HTML(http.StatusOK, "storage_edit.tmpl", data)
}

type StoragesPageArgs struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

type StoragesPostArgs struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

func StoragesEditPost(c *gin.Context) {
	var args model.Storage
	if err := c.ShouldBind(&args); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}

	fmt.Println(args)
	common.SuccessLayuiResp(c, 0, "ok")
}

func StoragesList(c *gin.Context) {
	var args StoragesPageArgs
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
