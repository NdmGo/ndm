package routers

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"ndm/internal/conf"
	"ndm/internal/handles"
	"ndm/internal/logs"
	"ndm/internal/middlewates"
	// "ndm/internal/utils"
	"ndm/internal/common"
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

	// Install Page
	g.POST("/check", handles.CheckConnectDb)
	g.GET("/install", handles.InstallPage)
	g.GET("/install_step1", handles.InstallStep1Page)
	g.POST("/install_step1", handles.PostInstallStep1Page)

	// Admin Page
	gnoauth := r.Group(conf.Http.SafePath, middlewates.SysIsInstalled, middlewates.PageNoAuth)
	gnoauth.GET("/login", handles.LoginPage)

	gauth := r.Group(conf.Http.SafePath, middlewates.PageAuth, middlewates.SysIsInstalled)
	gauth.GET("/", handles.IndexPage)
	gauth.GET("/storage", handles.StoragesPage)
	gauth.GET("/storage/edit", handles.StoragesEditPage)

	gauth.GET("/user", handles.UserPage)
	gauth.GET("/user/edit", handles.UserEditPage)

	// gauth.GET("/user", handles.StoragesPage)
	// user := g.Group("/user")
	// user.GET("/list", handles.ListUsers)

	gauth.GET("/setting", handles.SettingPage)
	gauth.GET("/task", handles.TaskPage)
	gauth.GET("/plugins", handles.PluginsPage)

}

func initRuoteApi(r *gin.Engine) {

	api := r.Group(conf.Http.ApiPath)
	api.Any("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	auth := api.Group("/auth")
	auth.POST("/login", handles.PostLogin)

	storage := api.Group("/storage")
	storage.GET("/list", handles.StoragesList)
	storage.POST("/edit", handles.CreateStorage)
	storage.POST("/create", handles.CreateStorage)
	storage.POST("/delete", handles.DeleteStorage)

	user := api.Group("/user")
	user.GET("/list", handles.ListUsers)
	user.POST("/create", handles.CreateUser)
	user.POST("/update", handles.UpdateUser)
	user.POST("/cancel_2fa", handles.Cancel2FAById)
	user.POST("/delete", handles.DeleteUser)
}

func InitRouters() {
	if !conf.Http.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	home := r.Group("", middlewates.SysIsInstalled)
	home.GET("/", func(c *gin.Context) {
		data := common.CommonVer()
		c.HTML(http.StatusOK, "home.tmpl", data)
	})

	r.SetTrustedProxies(nil)
	initAdminStaticPage(r)
	initRuoteApi(r)

	r.Run(fmt.Sprintf(":%d", conf.Http.Port))
}
