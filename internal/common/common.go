package common

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"ndm/internal/conf"
	"ndm/internal/utils"
)

type Resp[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type PageResp struct {
	Content interface{} `json:"content"`
	Total   int64       `json:"total"`
}

type LayuiResp[T any] struct {
	Code  int    `json:"code"`
	Count int64  `json:"count"`
	Msg   string `json:"msg"`
	Data  T      `json:"data"`
}

func CommonVer() map[string]interface{} {
	data := map[string]interface{}{
		"title":   "NDM存储管理",
		"version": conf.App.Version,
	}

	if !strings.EqualFold(conf.App.RunMode, "prod") {
		// 开发的时候开启
		data["version"] = fmt.Sprintf("%s_%s", conf.App.Version, utils.RandString(10))
	}

	data["admin_path"] = conf.Http.SafePath
	data["api_path"] = conf.Http.ApiPath
	return data
}

func ToJson(v interface{}) (d string) {
	rdata, _ := json.MarshalIndent(v, "", "  ")
	return string(rdata)
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
		Code: code,
		Msg:  hidePrivacy(err.Error()),
		Data: data,
	})
	c.Abort()
}

func ErrorStrResp(c *gin.Context, str string, code int, l ...bool) {
	if len(l) != 0 && l[0] {
		log.Error(str)
	}
	c.JSON(200, Resp[interface{}]{
		Code: code,
		Msg:  hidePrivacy(str),
		Data: nil,
	})
	c.Abort()
}

func SuccessResp(c *gin.Context, data ...interface{}) {
	if len(data) == 0 {
		c.JSON(200, Resp[interface{}]{
			Code: 200,
			Msg:  "success",
			Data: nil,
		})
		return
	}
	c.JSON(200, Resp[interface{}]{
		Code: 200,
		Msg:  "success",
		Data: data[0],
	})
}

func SuccessLayuiMsgResp(c *gin.Context, msg string, data ...interface{}) {
	if len(data) == 0 {
		c.JSON(200, LayuiResp[interface{}]{
			Code: 0,
			Msg:  msg,
			Data: nil,
		})
		return
	}

	c.JSON(200, LayuiResp[interface{}]{
		Code: 0,
		Msg:  msg,
		Data: data[0],
	})
}

func SuccessLayuiResp(c *gin.Context, count int64, msg string, data ...interface{}) {
	if len(data) == 0 {
		c.JSON(200, LayuiResp[interface{}]{
			Code:  0,
			Count: count,
			Msg:   msg,
			Data:  nil,
		})
		return
	}

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
