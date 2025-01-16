package handles

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ndm/internal/common"
	// "ndm/internal/conf"
)

func LoginPage(c *gin.Context) {
	data := common.CommonVer()
	c.HTML(http.StatusOK, "login.tmpl", data)
}
