package routers

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"ndm/internal/conf"
	"ndm/internal/utils"
)

func initAppRouter(e *gin.Engine) {
	if !utils.SliceContains([]string{"", "/"}, conf.Http.SafePath) {
		e.GET("/", func(c *gin.Context) {
			c.Redirect(302, conf.Http.SafePath)
		})
	}

	g := g.Group(conf.Http.SafePath)
	api := g.Group("/api")

	g.Any("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	api.Any("/pings", func(c *gin.Context) {
		c.String(200, "pong")
	})
}

func InitRouters() {
	engine := gin.Default()
	engine.SetTrustedProxies(nil)
	initAppRouter(engine)
	// gin.SetMode(gin.ReleaseMode)
	engine.Run(fmt.Sprintf(":%d", conf.Http.Port))
}
