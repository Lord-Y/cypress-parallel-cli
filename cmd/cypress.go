// Package cmd manage all commands required to launch cypress-parallel-cli
package cmd

import (
	"github.com/urfave/cli/v2"
)

// Cypress command options
func Cypress(c *cli.Context) (z *cli.Command) {
	return &cli.Command{
		Name:  "cypress",
		Usage: "options related to cypress command",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "api-url",
				Aliases:     []string{"a"},
				Value:       "http://127.0.0.1:8080",
				Usage:       "HTTP(s) api url of cypress-parallel-api",
				Destination: &cmd.ApiURL,
			},
			&cli.StringFlag{
				Name:        "repository",
				Aliases:     []string{"r"},
				Value:       "",
				Required:    true,
				Usage:       "HTTP(s) git repository (required)",
				Destination: &cmd.Repository,
			},
			&cli.StringFlag{
				Name:        "username",
				Aliases:     []string{"u"},
				Value:       "",
				Usage:       "Username to used to fetch repository if required",
				Destination: &cmd.Username,
			},
			&cli.StringFlag{
				Name:        "password",
				Aliases:     []string{"p"},
				Value:       "",
				Usage:       "Password to used to fetch repository if required",
				Destination: &cmd.Password,
			},
			&cli.StringFlag{
				Name:        "branch",
				Aliases:     []string{"b"},
				Value:       "",
				Usage:       "Branch in which specs are hold (required)",
				Destination: &cmd.Branch,
			},
			&cli.StringFlag{
				Name:        "specs",
				Aliases:     []string{"s"},
				Required:    true,
				Usage:       "Comma separated list of specs (required)",
				Destination: &cmd.Specs,
			},
			&cli.StringFlag{
				Name:        "uniq-id",
				Aliases:     []string{"uid"},
				Required:    true,
				Usage:       "Uniq ID to run cypress command (required)",
				Destination: &cmd.UniqID,
			},
			&cli.StringFlag{
				Name:        "browser",
				Aliases:     []string{"br"},
				Value:       "chrome",
				Usage:       "Default browser to use to run unit testing",
				Destination: &cmd.Browser,
			},
			&cli.StringFlag{
				Name:        "config-file",
				Aliases:     []string{"cf"},
				Value:       "cypress.json",
				Usage:       "Relative path of cypress config if not cypress.json",
				Destination: &cmd.ConfigFile,
			},
			&cli.BoolFlag{
				Name:        "report-back",
				Aliases:     []string{"rp"},
				Usage:       "Send result to api",
				Destination: &cmd.ReportBack,
			},
			&cli.IntFlag{
				Name:        "timeout",
				Aliases:     []string{"t"},
				Value:       10,
				Usage:       "Timeout after which the program will exit with error",
				Destination: &cmd.Timeout,
			},
		},
		Action: func(c *cli.Context) error {
			cmd.Run()
			return nil
		},
	}
}
