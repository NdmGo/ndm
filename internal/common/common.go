package common

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"ndm/internal/conf"
)

type Resp[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type PageResp struct {
	Content interface{} `json:"content"`
	Total   int64       `json:"total"`
}

func hidePrivacy(msg string) string {
	// for _, r := range conf.PrivacyReg {
	// 	msg = r.ReplaceAllStringFunc(msg, func(s string) string {
	// 		return strings.Repeat("*", len(s))
	// 	})
	// }
	return msg
}

func CommonVer() map[string]interface{} {
	data := map[string]interface{}{
		"title":   "NDM存储管理",
		"version": conf.App.Version,
	}
	return data
}

func ErrorStrResp(c *gin.Context, str string, code int, l ...bool) {
	if len(l) != 0 && l[0] {
		log.Error(str)
	}
	c.JSON(200, Resp[interface{}]{
		Code:    code,
		Message: hidePrivacy(str),
		Data:    nil,
	})
	c.Abort()
}

func SuccessResp(c *gin.Context, data ...interface{}) {
	if len(data) == 0 {
		c.JSON(200, Resp[interface{}]{
			Code:    200,
			Message: "success",
			Data:    nil,
		})
		return
	}
	c.JSON(200, Resp[interface{}]{
		Code:    200,
		Message: "success",
		Data:    data[0],
	})
}

func GetHttpReq(ctx context.Context) *http.Request {
	if c, ok := ctx.(*gin.Context); ok {
		return c.Request
	}
	return nil
}
