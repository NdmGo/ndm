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

func CommonVer() map[string]interface{} {
	data := map[string]interface{}{
		"title":   "NDM存储管理",
		"version": conf.App.Version,
	}

	data["admin_path"] = conf.Http.SafePath
	data["api_path"] = conf.Http.ApiPath
	return data
}

// ErrorResp is used to return error response
// @param l: if true, log error
func ErrorResp(c *gin.Context, err error, code int, l ...bool) {
	ErrorWithDataResp(c, err, code, nil, l...)
}

func hidePrivacy(msg string) string {
	// for _, r := range conf.PrivacyReg {
	// 	msg = r.ReplaceAllStringFunc(msg, func(s string) string {
	// 		return strings.Repeat("*", len(s))
	// 	})
	// }
	return msg
}

func ErrorWithDataResp(c *gin.Context, err error, code int, data interface{}, l ...bool) {
	if len(l) > 0 && l[0] {
		if conf.Http.Debug {
			log.Errorf("%+v", err)
		} else {
			log.Errorf("%v", err)
		}
	}
	c.JSON(200, Resp[interface{}]{
		Code:    code,
		Message: hidePrivacy(err.Error()),
		Data:    data,
	})
	c.Abort()
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

type LayuiResp[T any] struct {
	Code  int    `json:"code"`
	Count int64  `json:"count"`
	Msg   string `json:"msg"`
	Data  T      `json:"data"`
}

func SuccessLayuiResp(c *gin.Context, count int64, msg string, data ...interface{}) {
	c.JSON(200, LayuiResp[interface{}]{
		Code:  0,
		Count: count,
		Msg:   msg,
		Data:  data[0],
	})
}

func GetHttpReq(ctx context.Context) *http.Request {
	if c, ok := ctx.(*gin.Context); ok {
		return c.Request
	}
	return nil
}
