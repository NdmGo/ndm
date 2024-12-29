package routers

import (
	// "fmt"

	"github.com/gin-gonic/gin"
)

func InitRouters() {

	engine := gin.Default()
	engine.SetTrustedProxies(nil)
	gin.SetMode(gin.ReleaseMode)
	engine.Run(":5868")
}
