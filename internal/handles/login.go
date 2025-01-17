package handles

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"ndm/internal/common"
)

func LoginPage(c *gin.Context) {
	data := common.CommonVer()
	c.HTML(http.StatusOK, "login.tmpl", data)
}

func PostLogin(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")

	fmt.Println(username, password)

	data := common.CommonVer()
	c.HTML(http.StatusOK, "login.tmpl", data)
}
