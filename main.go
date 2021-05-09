package main

import (
	"os"

	"github.com/Lord-Y/cypress-parallel-cli/cmd"
	"github.com/Lord-Y/cypress-parallel-cli/logger"
	"github.com/rs/zerolog/log"

	cli "github.com/urfave/cli/v2"
)

// make vars public for unit testing
var (
	cypress     *cli.Command
	versionLong *cli.Command
)

func init() {
	logger.SetLoggerLogLevel()

	cypress = cmd.Cypress(&cli.Context{})
	versionLong = cmd.VersionLong(&cli.Context{})
}

func main() {
	app := cli.NewApp()
	app.Name = "cypress-parallel-cli"
	app.Usage = "A cli to start your e2e testing in parallel"
	app.Description = "cypress-parallel-cli enable user to reduce their e2e testing execution thanks to parallelization"
	app.Version = cmd.Version
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		cypress,
		versionLong,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Err(err).Msg("Error occured while executing the program")
	}
}
