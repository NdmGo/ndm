package handles

import (
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"ndm/internal/common"
	"ndm/internal/model"
	"ndm/internal/op"
)

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
	// OtpCode  string `json:"otp_code"`
}

func LoginPage(c *gin.Context) {
	data := common.CommonVer()
	c.HTML(http.StatusOK, "login.tmpl", data)
}

func PostLogin(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}

	req.Password = model.StaticHash(req.Password)
	loginHash(c, &req)
}

func loginHash(c *gin.Context, req *LoginReq) {
	// check username
	user, err := op.GetUserByName(req.Username)
	if err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	// validate password hash
	if err := user.ValidatePwdStaticHash(req.Password); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}

	// generate token
	token, err := common.GenerateToken(user)
	if err != nil {
		common.ErrorResp(c, err, 400, true)
		return
	}
	common.SuccessResp(c, gin.H{"token": token})
}
