package handles

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"ndm/internal/common"
	// "ndm/internal/model"
	// "ndm/internal/op"
)

func StoragesPage(c *gin.Context) {
	data := common.CommonVer()
	c.HTML(http.StatusOK, "storage.tmpl", data)
}

type StoragesArgs struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

func StoragesList(c *gin.Context) {
	var args StoragesArgs
	if err := c.ShouldBind(&args); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}

	fmt.Println(args)
}
