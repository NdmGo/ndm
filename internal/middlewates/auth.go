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

func PageAuth(c *gin.Context) {
	token, err := c.Cookie("token")

	url := fmt.Sprintf("%s/login", conf.Http.SafePath)
	if err != nil {

		c.Redirect(302, url)
		c.Next()
		return
	}

	userClaims, err := common.ParseToken(token)
	user, err := db.GetUserByName(userClaims.Username)
	if err != nil {
		common.ErrorResp(c, err, 401)
		c.Next()
		return
	}

	now_time := time.Now().Unix()
	token_expire := conf.Http.TokenExpiresIn

	time_expire := userClaims.PwdTS + (token_expire * 24 * 60 * 60)

	if time_expire > now_time {
		c.Redirect(302, url)
	}

	fmt.Println(time_expire, now_time)
	fmt.Println("user,", user)

	fmt.Println("token1,", userClaims.Username)
	fmt.Println("err,", err)

}
