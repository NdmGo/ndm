package routers

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"ndm/internal/common"
	"ndm/internal/conf"
	"ndm/internal/handles"
	"ndm/internal/logs"
	"ndm/internal/utils"
	"ndm/public"
)

// r *gin.RouterGroup
func ndmStatic(r *gin.Engine, noRoute func(handlers ...gin.HandlerFunc)) {
	folders := []string{"static"}
	r.Use(func(c *gin.Context) {
		for i := range folders {
			if strings.HasPrefix(c.Request.RequestURI, fmt.Sprintf("/%s/", folders[i])) {
				c.Header("Cache-Control", "public, max-age=15552000")
			}
		}
	})
	for i, folder := range folders {
		sub, err := fs.Sub(public.Static, folder)
		if err != nil {
			logs.Errorf("can't find folder: %s", folder)
		}
		r.StaticFS(fmt.Sprintf("/%s/", folders[i]), http.FS(sub))
	}

	noRoute(func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.Status(200)
		_, _ = c.Writer.WriteString("not find")
		c.Writer.Flush()
		c.Writer.WriteHeaderNow()
	})
}

func initStaticFunc(r *gin.Engine) {
	tmpl, err := template.ParseFS(public.Template, "template/default/**/*.tmpl", "template/default/*.tmpl")
	if err != nil {
		logs.Infof("load template: %s", err)
	}

	r.SetHTMLTemplate(tmpl)
	r.SetFuncMap(template.FuncMap{
		"safe": func(str string) template.HTML {
			return template.HTML(str)
		},
	})
}

func initAdminStaticPage(r *gin.Engine) {
	//static file
	ndmStatic(r, func(handlers ...gin.HandlerFunc) {
		r.NoRoute(handlers...)
	})

	initStaticFunc(r)
	g := r.Group(conf.Http.SafePath)
	// gauth := r.Group(conf.Http.SafePath, middlewates.Auth)

	g.GET("/", func(c *gin.Context) {
		if !conf.Security.InstallLock {
			c.Redirect(302, "/install")
		}

		data := common.CommonVer()
		c.HTML(http.StatusOK, "index.tmpl", data)
	})

	g.POST("/check", handles.CheckConnectDb)
	g.GET("/install", handles.InstallPage)
	g.GET("/install_step1", handles.InstallStep1Page)
	g.POST("/install_step1", handles.PostInstallStep1Page)
	g.GET("/login", handles.InstallPage)
}

func InitRouters() {
	if !conf.Http.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.SetTrustedProxies(nil)
	if !utils.SliceContains([]string{"", "/"}, conf.Http.SafePath) {
		// r.GET("/", func(c *gin.Context) {
		// 	c.Redirect(302, conf.Http.SafePath)
		// })
	}
	initAdminStaticPage(r)

	// api := r.Group("/api")
	// api.Any("/ping", func(c *gin.Context) {
	// 	c.String(200, "pong")
	// })

	// api.Any("/pings", func(c *gin.Context) {
	// 	c.String(200, "pong")
	// })

	r.Run(fmt.Sprintf(":%d", conf.Http.Port))
}
