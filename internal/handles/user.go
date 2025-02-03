package handles

import (
	// "fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ndm/internal/common"
	"ndm/internal/db"
	"ndm/internal/model"
	// "ndm/internal/op"
)

func UserPage(c *gin.Context) {
	data := common.CommonVer()
	c.HTML(http.StatusOK, "user.tmpl", data)
}

func UserEditPage(c *gin.Context) {
	data := common.CommonVer()
	c.HTML(http.StatusOK, "user_edit.tmpl", data)
}

type UserArgs struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

func ListUsers(c *gin.Context) {
	var args model.PageReq
	if err := c.ShouldBind(&args); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}

	args.Validate()
	users, total, err := db.GetUsers(args.Page, args.Size)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessLayuiResp(c, total, "ok", users)
}

func CreateUser(c *gin.Context) {
	var req model.User
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if req.IsAdmin() || req.IsGuest() {
		common.ErrorStrResp(c, "admin or guest user can not be created", 400, true)
		return
	}
	req.SetPassword(req.Password)
	req.Password = ""
	req.Authn = "[]"
	if err := db.CreateUser(&req); err != nil {
		common.ErrorResp(c, err, 500, true)
	} else {
		common.SuccessResp(c)
	}
}

func DeleteUser(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if err := db.DeleteUserById(int64(id)); err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessResp(c)
}
