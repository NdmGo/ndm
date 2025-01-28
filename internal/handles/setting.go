package handles

import (
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"ndm/internal/common"
	// "ndm/internal/db"
	// "ndm/internal/model"
	// "ndm/internal/op"
)

func SettingPage(c *gin.Context) {
	data := common.CommonVer()
	c.HTML(http.StatusOK, "setting.tmpl", data)
}
