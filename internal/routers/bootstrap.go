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
	"ndm/internal/middlewares"
	"ndm/internal/sign"
	"ndm/internal/stream"
	// "ndm/internal/utils"
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

	noRoute(handles.HomePage)
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

	initStaticFunc(r)
	g := r.Group(conf.Http.SafePath)

	// Install Page
	g.POST("/check", handles.CheckConnectDb)
	g.GET("/install", handles.InstallPage)
	g.GET("/install_step1", handles.InstallStep1Page)
	g.POST("/install_step1", handles.PostInstallStep1Page)

	// Admin Page
	gnoauth := r.Group(conf.Http.SafePath, middlewares.SysIsInstalled, middlewares.PageNoAuth)
	gnoauth.GET("/login", handles.LoginPage)

	gauth := r.Group(conf.Http.SafePath, middlewares.PageAuth, middlewares.SysIsInstalled)
	gauth.GET("/", handles.AdminPage)
	gauth.GET("/storage", handles.StoragesPage)
	gauth.GET("/storage/create", handles.StoragesCreatePage)
	gauth.GET("/storage/edit", handles.StoragesEditPage)

	gauth.GET("/user", handles.UserPage)
	gauth.GET("/user/edit", handles.UserEditPage)

	gauth.GET("/setting", handles.SettingPage)
	gauth.GET("/setting/:action", handles.SettingPage)

	gauth.GET("/logs", handles.LogsPage)

	gauth.GET("/tasks", handles.TasksPage)
	gauth.GET("/tasks/create", handles.TasksCreatePage)
	gauth.GET("/tasks/edit", handles.TasksEditPage)
	gauth.GET("/tasks/:action", handles.TasksPage)
	gauth.GET("/plugins", handles.PluginsPage)

	//static file
	ndmStatic(r, func(handlers ...gin.HandlerFunc) {
		r.NoRoute(handlers...)
	})

}

func initFs(fs *gin.RouterGroup) {
	fs.Any("/list", handles.FsList)
	fs.Any("/get", handles.FsGet)
}

func initRuoteApi(r *gin.Engine) {

	g := r.Group(conf.Http.ApiPath)
	auth := g.Group("", middlewares.Auth)

	g.Any("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	login := g.Group("/auth")
	login.POST("/login", handles.PostLogin)

	user := g.Group("/user")
	user.GET("/list", handles.ListUsers)
	user.POST("/create", handles.CreateUser)
	user.POST("/update", handles.UpdateUser)
	user.POST("/cancel_2fa", handles.Cancel2FAById)
	user.POST("/delete", handles.DeleteUser)

	storage := g.Group("/storage")
	storage.GET("/list", handles.StoragesList)
	storage.POST("/update", handles.UpdateStorage)
	storage.POST("/create", handles.CreateStorage)
	storage.POST("/delete", handles.DeleteStorage)
	storage.POST("/trigger_disable", handles.TriggerDisabledStorage)

	logs := g.Group("/logs")
	logs.GET("/list", handles.LogsList)
	logs.POST("/delete", handles.DeleteLogs)
	logs.POST("/truncate", handles.TruncateLogs)
	logs.POST("/get", handles.GetLogs)

	tasks := g.Group("/tasks")
	tasks.GET("/list", handles.TasksList)
	tasks.POST("/create", handles.CreateTasks)
	tasks.POST("/delete", handles.DeleteTasks)
	tasks.POST("/done", handles.DoneTasks)

	downloadLimiter := middlewares.DownloadRateLimiter(stream.ClientDownloadLimit)
	signCheck := middlewares.Down(sign.Verify)
	r.GET("/d/*path", signCheck, downloadLimiter, handles.Down)
	r.GET("/p/*path", signCheck, downloadLimiter, handles.Proxy)
	r.HEAD("/d/*path", signCheck, handles.Down)
	r.HEAD("/p/*path", signCheck, handles.Proxy)

	initFs(auth.Group("/fs"))
}

func InitRouters() {
	if !conf.Http.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	home := r.Group("", middlewares.SysIsInstalled)
	home.GET("/", handles.HomePage)

	r.SetTrustedProxies(nil)

	initRuoteApi(r)
	initAdminStaticPage(r)

	r.Run(fmt.Sprintf(":%d", conf.Http.Port))
}
