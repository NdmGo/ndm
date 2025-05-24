package main

import (
	"log"
	"os"

	"github.com/urfave/cli"

	"ndm/internal/cmd"
	"ndm/internal/conf"
)

const Version = "0.0.1.28.7"
const AppName = "ndm"

func init() {
	conf.App.Version = Version
	conf.App.Name = AppName
}

func main() {

	app := cli.NewApp()
	app.Name = AppName
	app.Version = Version
	app.Usage = "ndm service"
	app.Commands = []cli.Command{
		cmd.Web,
	}

	if err := app.Run(os.Args); err != nil {
		log.Println("Failed to start application: ", err)
	}

}
