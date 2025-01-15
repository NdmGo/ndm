package middlewates

import (
	"fmt"
	// "crypto/subtle"

	"github.com/gin-gonic/gin"
	// "ndm/internal/conf"
	// "ndm/internal/logs"
	// "ndm/internal/model"
	// "ndm/internal/op"
)

func Auth(c *gin.Context) {
	token := c.GetHeader("Authorization")

	fmt.Println("token", token)

}
