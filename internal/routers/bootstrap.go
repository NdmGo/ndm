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

func initStaticPage(g *gin.RouterGroup) {
	g.GET("/", func(c *gin.Context) {
		if !conf.Security.InstallLock {
			c.Redirect(302, "/install")
		}

		data := common.CommonVer()
		c.HTML(http.StatusOK, "index.tmpl", data)
	})

	g.GET("/install", func(c *gin.Context) {
		data := common.CommonVer()
		c.HTML(http.StatusOK, "install.tmpl", data)
	})

	g.POST("/install_step1", func(c *gin.Context) {
		username := c.PostForm("username")

		fmt.Println(username)

		// c.JSON(http.StatusOK, gin.H{
		// 	"username": username,
		// 	"password": password,
		// 	"age":      age,
		// })

		common.SuccessResp(c, gin.H{"token": "安装成功!"})
	})

	g.GET("/install_step1", func(c *gin.Context) {
		data := common.CommonVer()
		c.HTML(http.StatusOK, "install_step1.tmpl", data)
	})
}

func InitRouters() {
	if !conf.Http.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	initStaticFunc(r)

	r.SetTrustedProxies(nil)
	if !utils.SliceContains([]string{"", "/"}, conf.Http.SafePath) {
		r.GET("/", func(c *gin.Context) {
			c.Redirect(302, conf.Http.SafePath)
		})
	}
	g := r.Group(conf.Http.SafePath)
	initStaticPage(g)

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
