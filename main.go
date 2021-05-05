package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Lord-Y/cypress-parallel-cli/cmd"
	"github.com/Lord-Y/cypress-parallel-cli/logger"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	cli "github.com/urfave/cli/v2"
)

// make vars public for unit testing
var (
	cypress *cli.Command
)

func init() {
	logger.SetLoggerLogLevel()

	output := zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %s |", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	log.Logger = zerolog.New(output).With().Timestamp().Logger()

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
		log.Fatal().Err(err).Msg("Error occured while executing the program")
	}
}
