package main

import (
	"log"
	"os"

	"github.com/Lord-Y/cypress-parallel-cli/cmd"

	cli "github.com/urfave/cli/v2"
)

// make vars public for unit testing
var (
	cypress *cli.Command
)

func init() {
	cypress = cmd.Cypress(&cli.Context{})
}

func main() {
	app := cli.NewApp()
	app.Name = "cypress-parallel-cli"
	app.Usage = "A cli to start your e2e testing in parallel"
	app.Description = "cypress-parallel-cli enable user to reduce their e2e testing execution thanks to parallelization"
	app.Version = "0.0.1"
	app.EnableBashCompletion = true

	app.Commands = []*cli.Command{
		cypress,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
