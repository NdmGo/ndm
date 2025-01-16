package middlewates

import (
	"fmt"
	// "crypto/subtle"

	"github.com/gin-gonic/gin"
	"ndm/internal/conf"
	// "ndm/internal/logs"
	// "ndm/internal/model"
	// "ndm/internal/op"
)

func PageAuth(c *gin.Context) {
	cookie, err := c.Cookie("token")

	if err != nil {
		url := fmt.Sprintf("%s/login", conf.Http.SafePath)
		c.Redirect(302, url)
		c.Next()
		return
	}

	fmt.Println("cookie,", cookie)

}
