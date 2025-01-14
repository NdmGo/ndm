package middlewates

import (
	// "crypto/subtle"

	"github.com/gin-gonic/gin"

	"ndm/internal/conf"
	"ndm/internal/logs"
	"ndm/internal/model"
	"ndm/internal/op"
	"ndm/internal/setting"
)

func Auth(c *gin.Context) {
	token := c.GetHeader("Authorization")

	fmt.Println("token", token)

}
