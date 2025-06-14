package middlewares

import (
	"fmt"
	"time"
	// "crypto/subtle"

	"github.com/gin-gonic/gin"
	"ndm/internal/common"
	"ndm/internal/conf"
	"ndm/internal/utils"
	// "ndm/internal/logs"
	"ndm/internal/db"
	// "ndm/internal/model"
	// "ndm/internal/op"

	log "github.com/sirupsen/logrus"
)

func Auth(c *gin.Context) {
	token := c.GetHeader("Authorization")

	userClaims, err := common.ParseToken(token)
	if err != nil {
		common.ErrorResp(c, err, 401)
		c.Abort()
		return
	}

	user, err := db.GetUserByName(userClaims.Username)
	if err != nil {
		common.ErrorResp(c, err, 401)
		c.Abort()
		return
	}

	// validate password timestamp
	if userClaims.PwdTS != user.PwdTS {
		common.ErrorStrResp(c, "Password has been changed, login please", 401)
		c.Abort()
		return
	}
	if user.Disabled {
		common.ErrorStrResp(c, "Current user is disabled, replace please", 401)
		c.Abort()
		return
	}

	now_time := time.Now().Unix()
	token_expire_time := userClaims.RegisteredClaims.ExpiresAt.Unix()
	if now_time > token_expire_time {
		common.ErrorStrResp(c, "Login has expired, login please", 401)
		c.Abort()
		return
	}

	c.Set("user", user)
	log.Debugf("use login token: %+v", user)
	c.Next()
}

func PageNoAuth(c *gin.Context) {
	url := fmt.Sprintf("%s", conf.Http.SafePath)

	token, err := c.Cookie("token")
	if err != nil {
		c.Next()
		return
	}

	userClaims, err := common.ParseToken(token)
	if err != nil {
		c.Next()
		return
	}

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
	if err != nil {
		c.Redirect(302, url)
		c.Next()
		return
	}

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

func SysIsInstalled(c *gin.Context) {
	conf_path := conf.WorkDir()
	custom_dir := fmt.Sprintf("%s/custom", conf_path)
	install_url := fmt.Sprintf("%s/install", conf.Http.SafePath)
	if !utils.IsExist(custom_dir) {
		c.Redirect(302, install_url)
		return
	}
}
