package routers

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"ndm/internal/conf"
	"ndm/internal/logs"
	"ndm/internal/utils"
	"ndm/public"
)

func ndmStatic(r *gin.RouterGroup, noRoute func(handlers ...gin.HandlerFunc)) {
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

func InitRouters() {
	r := gin.Default()

	tmpl, err := template.ParseFS(public.Template, "template/default/*")
	if err != nil {
		logs.Infof("load template err: %s", err)
	}

	// r.Delims("{[", "]}")
	r.SetHTMLTemplate(tmpl)
	r.SetFuncMap(template.FuncMap{
		"title": "测试",
		"AppName": func() string {
			return conf.App.Name
		},
		"AppVer": func() string {
			return conf.App.Version
		},
		"LoadTimes": func(startTime time.Time) string {
			return fmt.Sprint(time.Since(startTime).Nanoseconds()/1e6) + "ms"
		},
		"Join": strings.Join,
		"DateFmtLong": func(t time.Time) string {
			return t.Format(time.RFC1123Z)
		},
		"DateFmtShort": func(t time.Time) string {
			return t.Format("Jan 02, 2006")
		},
	})
	r.SetTrustedProxies(nil)
	if !utils.SliceContains([]string{"", "/"}, conf.Http.SafePath) {
		r.GET("/", func(c *gin.Context) {
			c.Redirect(302, conf.Http.SafePath)
		})
	}
	g := r.Group(conf.Http.SafePath)
	g.Any("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})

	// g.Any("/install", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "install.tmpl", gin.H{})
	// })

	if !conf.Http.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// api := g.Group("/api")

	// g.Any("/ping", func(c *gin.Context) {
	// 	c.String(200, "pong")
	// })

	// api.Any("/pings", func(c *gin.Context) {
	// 	c.String(200, "pong")
	// })

	ndmStatic(g, func(handlers ...gin.HandlerFunc) {
		r.NoRoute(handlers...)
	})

	r.Run(fmt.Sprintf(":%d", conf.Http.Port))
}
