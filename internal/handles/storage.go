package handles

import (
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"ndm/internal/common"
	// "ndm/internal/model"
	// "ndm/internal/op"
)

func StoragePage(c *gin.Context) {
	data := common.CommonVer()
	c.HTML(http.StatusOK, "storage.tmpl", data)
}

func StorageList(page, size int64) {

}
