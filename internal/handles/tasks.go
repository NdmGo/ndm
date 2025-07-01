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
	"ndm/internal/errs"
	"ndm/internal/model"
	"ndm/internal/op"
)

func TasksPage(c *gin.Context) {
	data := common.CommonVer()

	action := c.Param("action")
	data["task_page"] = action
	c.HTML(http.StatusOK, "tasks.tmpl", data)
}

func TasksCreatePage(c *gin.Context) {
	data := common.CommonVer()
	c.HTML(http.StatusOK, "tasks_create.tmpl", data)
}

func TasksList(c *gin.Context) {
	var args model.PageReq
	if err := c.ShouldBind(&args); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}

	data_list, total, err := db.GetTasksList(args.Page, args.Size)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessLayuiResp(c, total, "ok", data_list)
}

func CreateTasks(c *gin.Context) {
	var req model.Tasks
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

	_, err := db.GetTasksByMountPath(req.MountPath)
	if err == nil {
		common.ErrorWithDataResp(c, errors.New("任务已经存在!"), 500, gin.H{}, true)
		return
	}

	if id, err := op.CreateTasks(req); err != nil {
		common.ErrorWithDataResp(c, err, 500, gin.H{"id": id}, true)
	} else {
		common.SuccessResp(c, gin.H{"id": id})
	}
}

func DoneTasks(c *gin.Context) {
	var req model.Tasks
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}

	if strings.EqualFold(req.MountPath, "") {
		common.ErrorWithDataResp(c, errs.MountPathCannotEmpty, 500, gin.H{}, true)
		return
	}

	_, err := db.GetTasksByMountPath(req.MountPath)
	if err != nil {
		common.ErrorWithDataResp(c, err, 500, gin.H{}, true)
		return
	}

	err = op.DoneTasksBackup(c, req.MountPath)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessResp(c, "ok")
}

func DeleteTasks(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if err := op.DeleteTasksById(int64(id)); err != nil {
		common.ErrorResp(c, err, 500, true)
		return
	}
	common.SuccessResp(c)
}
