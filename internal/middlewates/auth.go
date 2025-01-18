package middlewates

import (
	"fmt"
	"time"
	// "crypto/subtle"

	"github.com/gin-gonic/gin"
	"ndm/internal/common"
	"ndm/internal/conf"
	// "ndm/internal/logs"
	"ndm/internal/db"
	// "ndm/internal/model"
	// "ndm/internal/op"
)

func PageNoAuth(c *gin.Context) {
	url := fmt.Sprintf("%s", conf.Http.SafePath)

	token, err := c.Cookie("token")
	if err != nil {
		c.Next()
		return
	}

	userClaims, err := common.ParseToken(token)
	// fmt.Println(userClaims, err)
	_, err = db.GetUserByName(userClaims.Username)
	if err != nil {
		c.Next()
		return
	}

	now_time := time.Now().Unix()
	token_expire_time := userClaims.RegisteredClaims.ExpiresAt.Unix()
	if now_time < token_expire_time {
		c.Redirect(302, url)
		c.Next()
		return
	}
}

func PageAuth(c *gin.Context) {
	url := fmt.Sprintf("%s/login", conf.Http.SafePath)

	token, err := c.Cookie("token")
	if err != nil {
		c.Redirect(302, url)
		c.Next()
		return
	}

	userClaims, err := common.ParseToken(token)
	_, err = db.GetUserByName(userClaims.Username)
	if err != nil {
		c.Redirect(302, url)
		c.Next()
		return
	}

	now_time := time.Now().Unix()
	token_expire_time := userClaims.RegisteredClaims.ExpiresAt.Unix()
	if now_time > token_expire_time {
		c.Redirect(302, url)
		c.Next()
		return
	}
}
