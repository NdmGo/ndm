package middlewares

import (
	"strings"

	"ndm/internal/conf"
	"ndm/internal/setting"

	"ndm/internal/common"
	// "ndm/internal/errs"
	"ndm/internal/model"
	"ndm/internal/op"
	"ndm/pkg/utils"

	"github.com/gin-gonic/gin"
	// "github.com/pkg/errors"
)

func Down(verifyFunc func(string, string) error) func(c *gin.Context) {
	return func(c *gin.Context) {
		rawPath := parsePath(c.Param("path"))
		c.Set("path", rawPath)
		// verify sign
		if needSign(rawPath) {
			s := c.Query("sign")
			err := verifyFunc(rawPath, strings.TrimSuffix(s, "/"))
			if err != nil {
				common.ErrorResp(c, err, 401)
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// TODO: implement
// path maybe contains # ? etc.
func parsePath(path string) string {
	return utils.FixAndCleanPath(path)
}

func needSign(path string) bool {
	if setting.GetBool(conf.SignAll) {
		return true
	}
	if op.IsStorageSignEnabled(path) {
		return true
	}
	return true
}

func needSignBak(meta *model.Meta, path string) bool {
	if setting.GetBool(conf.SignAll) {
		return true
	}
	if op.IsStorageSignEnabled(path) {
		return true
	}
	if meta == nil || meta.Password == "" {
		return false
	}
	if !meta.PSub && path != meta.Path {
		return false
	}
	return true
}
