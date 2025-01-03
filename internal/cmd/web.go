package cmd

import (
	"github.com/urfave/cli"

	"ndm/internal/conf"
	"ndm/internal/routers"
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
	conf.InitConf("")
	routers.InitRouters()
	return nil
}
