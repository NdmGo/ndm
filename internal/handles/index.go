package handles

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ndm/internal/common"
	"ndm/internal/conf"
)

func IndexPage(c *gin.Context) {
	if !conf.Security.InstallLock {
		c.Redirect(302, "/install")
	}
	data := common.CommonVer()
	c.HTML(http.StatusOK, "index.tmpl", data)
}
