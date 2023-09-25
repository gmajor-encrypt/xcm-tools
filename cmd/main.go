package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := setup()
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func setup() *cli.App {
	app := cli.NewApp()
	app.Name = "Xcm tools"
	app.Usage = "Xcm tools"
	app.Action = func(*cli.Context) error { return nil }
	app.Flags = []cli.Flag{}
	app.Commands = subCommands()
	return app
}
