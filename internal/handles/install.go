package handles

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ndm/internal/common"
	"ndm/internal/conf"
	"ndm/internal/db"
	userdata "ndm/internal/routers/data"
)

func InstallPage(c *gin.Context) {
	if conf.Security.InstallLock {
		c.Redirect(302, "/")
	}
	data := common.CommonVer()
	c.HTML(http.StatusOK, "install.tmpl", data)
}

func InstallStep1Page(c *gin.Context) {
	if conf.Security.InstallLock {
		c.Redirect(302, "/")
	}

	data := common.CommonVer()
	c.HTML(http.StatusOK, "install_step1.tmpl", data)
}

func PostInstallStep1Page(c *gin.Context) {
	install_data := make(map[string]string, 0)
	install_data["type"] = c.PostForm("type")
	install_data["hostname"] = c.PostForm("hostname")
	install_data["hostport"] = c.PostForm("hostport")
	install_data["dbname"] = c.PostForm("dbname")
	install_data["username"] = c.PostForm("username")
	install_data["password"] = c.PostForm("password")
	install_data["table_prefix"] = c.PostForm("table_prefix")

	init_account := c.PostForm("account")
	init_pass := c.PostForm("pass")
	err := conf.InstallConf(install_data)
	if err != nil {
		common.ErrorStrResp(c, err.Error(), -1)
		return
	}

	if conf.Security.InstallLock {
		db.InitDb()
		userdata.InitAdmin(init_account, init_pass)
	}

	common.SuccessResp(c, gin.H{"token": "安装成功!"})
}

func CheckConnectDb(c *gin.Context) {
	data := make(map[string]string, 0)
	data["type"] = c.PostForm("type")
	data["hostname"] = c.PostForm("hostname")
	data["hostport"] = c.PostForm("hostport")
	data["dbname"] = c.PostForm("dbname")
	data["username"] = c.PostForm("username")
	data["password"] = c.PostForm("password")
	data["dbprefix"] = c.PostForm("dbprefix")
	err := db.CheckDbConnnect(data)
	if err != nil {
		common.ErrorStrResp(c, err.Error(), -1)
		return
	}
	common.SuccessResp(c)
}
