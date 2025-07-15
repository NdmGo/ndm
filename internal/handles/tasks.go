package handles

import (
	// "errors"
	// "fmt"
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

func TasksEditPage(c *gin.Context) {
	data := common.CommonVer()

	tidStr := c.Query("id")
	tid, err := strconv.Atoi(tidStr)
	if err != nil {
		tid = 0
	}
	task, err := db.GetTasksById(int64(tid))
	if err == nil {
		data["task"] = task
	}
	c.HTML(http.StatusOK, "tasks_edit.tmpl", data)
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

	if req.MpId == 0 {
		common.ErrorWithDataResp(c, errs.MountPathCannotEmpty, 500, gin.H{
			"id": 0,
		}, true)
		return
	}

	_, err := db.GetTasksByMpId(req.MpId)
	if err == nil {
		common.ErrorWithDataResp(c, errs.TaskAlredyExists, 500, gin.H{}, true)
		return
	}

	if id, err := op.CreateTasks(req); err != nil {
		common.ErrorWithDataResp(c, err, 500, gin.H{"id": id}, true)
	} else {
		common.SuccessResp(c, gin.H{"id": id})
	}
}

func UpdateTask(c *gin.Context) {
	var req model.Tasks
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}

	if err := db.UpdateTask(&req); err != nil {
		common.ErrorResp(c, err, 500)
	} else {
		common.SuccessResp(c)
	}
}

func DoneTasks(c *gin.Context) {
	var req model.Tasks
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}

	if req.MpId == 0 {
		common.ErrorWithDataResp(c, errs.MountPathCannotEmpty, 500, gin.H{}, true)
		return
	}

	_, err := db.GetTasksByMpId(req.MpId)
	if err != nil {
		common.ErrorWithDataResp(c, err, 500, gin.H{}, true)
		return
	}

	storage, err := db.GetStorageById(req.MpId)
	if err != nil {
		common.ErrorWithDataResp(c, err, 500, gin.H{}, true)
		return
	}

	if strings.EqualFold(storage.Driver, "local") {
		err = op.DoneTasksSync(c, storage.MountPath)
	} else {
		err = op.DoneTasksBackup(c, storage.MountPath)
	}

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
