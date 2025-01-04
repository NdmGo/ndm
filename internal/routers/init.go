package routers

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"ndm/internal/conf"
)

func InitRouters() {
	engine := gin.Default()
	engine.SetTrustedProxies(nil)
	// gin.SetMode(gin.ReleaseMode)
	port := fmt.Sprintf(":%d", conf.Http.Port)
	engine.Run(port)
}
