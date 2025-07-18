package cmd

import (
	"net/http"

	"ndm/drivers"
	"ndm/internal/conf"
	"ndm/internal/crontab"
	"ndm/internal/db"
	"ndm/internal/logs"
	"ndm/internal/routers"
	userdata "ndm/internal/routers/data"

	"github.com/urfave/cli"
)

var Web = cli.Command{
	Name:        "web",
	Usage:       "This command start web service",
	Description: `Start web service`,
	Action:      runWeb,
	Flags: []cli.Flag{
		stringFlag("config, c", "", "custom configuration file path"),
	},
}

func runWeb(c *cli.Context) error {
	conf.InitConf(c.String("config"))
	logs.InitLog()
	drivers.All()

	if conf.Security.InstallLock {
		db.InitDb()
		userdata.InitAdmin("admin", "admin")
	}

	if conf.App.RunMode != "prod" {
		go func() {
			http.ListenAndServe("localhost:6060", nil)
		}()
	}

	routers.LoadStorages()
	crontab.Load()
	routers.InitRouters()

	return nil
}
