package handles

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"ndm/internal/common"
	"ndm/internal/db"
	"ndm/internal/model"
	"ndm/internal/op"
)

func StoragesPage(c *gin.Context) {
	data := common.CommonVer()
	c.HTML(http.StatusOK, "storage.tmpl", data)
}

func StoragesCreatePage(c *gin.Context) {
	data := common.CommonVer()
	c.HTML(http.StatusOK, "storage_create.tmpl", data)
}

func StoragesEditPage(c *gin.Context) {
	data := common.CommonVer()
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		id = 0
	}
	storage, err := db.GetStorageById(int64(id))
	if err == nil {
		data["storage"] = storage
	}

	tpl_name := fmt.Sprintf("storage_edit_%s.tmpl", storage.Driver)
	c.HTML(http.StatusOK, tpl_name, data)
}

func StoragesEditPost(c *gin.Context) {
	data := common.CommonVer()
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		id = 0
	}
	storage, err := db.GetStorageById(int64(id))
	if err == nil {
		data["storage"] = storage
	}
	common.SuccessLayuiResp(c, 0, "ok")
}

func CreateStorage(c *gin.Context) {
	var req model.Storage
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}

	if strings.EqualFold(req.MountPath, "") {
		common.ErrorWithDataResp(c, errors.New("挂载路径不能为空!"), 500, gin.H{
			"id": 0,
		}, true)
		return
	}

	if id, err := op.CreateStorage(c, req); err != nil {
		common.ErrorWithDataResp(c, err, 500, gin.H{
			"id": id,
		}, true)
	} else {
		common.SuccessResp(c, gin.H{
			"id": id,
		})
	}
}

func UpdateStorage(c *gin.Context) {

}

func DeleteStorage(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if err := op.DeleteStorageById(c, int64(id)); err != nil {
		common.ErrorResp(c, err, 500, true)
		return
	}
	common.SuccessResp(c)
}

func TriggerDisabledStorage(c *gin.Context) {
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if err := op.TriggerDisabledStorageById(c, int64(id)); err != nil {
		common.ErrorResp(c, err, 500, true)
		return
	}
	common.SuccessResp(c)
}

func StoragesList(c *gin.Context) {
	var args model.PageReq
	if err := c.ShouldBind(&args); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}

	// fmt.Println(args)
	storages, total, err := db.GetStorages(args.Page, args.Size)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessLayuiResp(c, total, "ok", storages)
}
