package routers

import (
	// "fmt"

	"github.com/gin-gonic/gin"
)

func InitRouters() {

	r := gin.Default()

	r.Run(":5868")
}
