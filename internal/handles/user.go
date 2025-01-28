package handles

import (
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"ndm/internal/common"
	"ndm/internal/db"
	// "ndm/internal/model"
	// "ndm/internal/op"
)

func UserPage(c *gin.Context) {
	data := common.CommonVer()
	c.HTML(http.StatusOK, "user.tmpl", data)
}

type UserArgs struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

func UserList(c *gin.Context) {
	var args UserArgs
	if err := c.ShouldBind(&args); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}

	storages, total, err := db.GetStorages(args.Page, args.Size)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	common.SuccessResp(c, common.PageResp{
		Content: storages,
		Total:   total,
	})
}
