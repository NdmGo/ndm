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
	idStr := c.Query("id")
	uid, err := strconv.Atoi(idStr)
	if err != nil {
		uid = 0
	}
	user, err := db.GetUserById(int64(uid))
	if err == nil {
		data["user"] = user
	}
	c.HTML(http.StatusOK, "user_edit.tmpl", data)
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

func UpdateUser(c *gin.Context) {
	var req model.User
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	user, err := db.GetUserById(req.ID)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	if user.Role != req.Role {
		common.ErrorStrResp(c, "role can not be changed", 400)
		return
	}
	if req.Password == "" {
		req.PwdHash = user.PwdHash
		req.Salt = user.Salt
	} else {
		req.SetPassword(req.Password)
		req.Password = ""
	}
	if req.OtpSecret == "" {
		req.OtpSecret = user.OtpSecret
	}
	if req.Disabled && req.IsAdmin() {
		common.ErrorStrResp(c, "admin user can not be disabled", 400)
		return
	}
	if err := db.UpdateUser(&req); err != nil {
		common.ErrorResp(c, err, 500)
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
	common.SuccessLayuiMsgResp(c, "user deleted successfully")
}
