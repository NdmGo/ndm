package handles

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"ndm/internal/common"
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

	fmt.Println(req)

	username := c.PostForm("username")
	password := c.PostForm("password")

	fmt.Println(username, password)

	data := common.CommonVer()
	c.HTML(http.StatusOK, "login.tmpl", data)
}
