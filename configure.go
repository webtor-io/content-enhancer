package main

import (
	"github.com/urfave/cli"
)

func configure(app *cli.App) {
	runCmd := makeRunCMD()
	app.Commands = []cli.Command{runCmd}
}
