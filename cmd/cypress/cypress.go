// Package cypress assemble all commands required to run cypress unit testing
package cypress

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Lord-Y/cypress-parallel-cli/git"
	"github.com/Lord-Y/cypress-parallel-cli/logger"
	"github.com/rs/zerolog/log"
)

// Cypress requirements to run cypress command
type Cypress struct {
	ApiURL     string // HTTP(s) api url of cypress-parallel-api
	Repository string // HTTP(s) git repository
	Username   string // Username to use to fetch repository if required
	Password   string // Password to use to fetch repository if required
	Branch     string // Branch in which specs are hold
	Specs      string // Comma separated list of specs
	UniqID     string // Uniq ID to run cypress command
	Browser    string // Default browser to use to run unit testing
	ConfigFile string // Relative path of cypress config if not cypress.json
	ReportBack bool   // Notify api with cypress results
}

func init() {
	logger.SetLoggerLogLevel()
}

// Run will run cypress command
func (c *Cypress) Run() {
	var (
		gc git.Repository
	)

	gc.Repository = c.Repository
	gc.Username = c.Username
	gc.Password = c.Password
	gc.Branch = c.Branch

	gitdir, err := gc.Clone()
	if err != nil {
		log.Fatal().Err(err).Msg("Error occured while cloning git repository")
		return
	}
	defer os.RemoveAll(gitdir)
	log.Debug().Msgf("Git temp dir %s", gitdir)

	err = os.Chdir(gitdir)
	if err != nil {
		log.Fatal().Err(err).Msg("Error occured while chdir git repository")
		return
	}

	if c.ConfigFile != "" {
		var info os.FileInfo
		if info, err = os.Stat(fmt.Sprintf("%s/%s", gitdir, c.ConfigFile)); os.IsNotExist(err) {
			log.Fatal().Err(err).Msgf("Error occured while checking config file %s", c.ConfigFile)
			return
		}
		if !info.Mode().IsRegular() {
			log.Fatal().Err(err).Msgf("Error occured while checking config file %s", c.ConfigFile)
			return
		}
	}

	execUninstallCmd := exec.Command(
		"npm",
		"uninstall",
		"cypress",
		"&&",
		"npm",
		"uninstall",
		"-D",
		"cypress",
	)

	if err := execUninstallCmd.Run(); err != nil {
		log.Fatal().Err(err).Msg("Error occured while forcing uninstall of local cypress package")
		return
	}

	output, err := exec.Command(
		"npm",
		"install",
	).Output()
	log.Debug().Msgf("NPM install output %s", string(output))

	if err != nil {
		log.Fatal().Err(err).Msgf("Error occured while installing user packages")
		return
	}

	output, err = exec.Command(
		"npm",
		"install",
		"--save-dev",
		"mochawesome",
	).Output()
	log.Debug().Msgf("Mochawesome install output %s", string(output))

	if err != nil {
		log.Fatal().Err(err).Msgf("Error occured while installing mochawesome")
		return
	}

	specs := strings.Split(c.Specs, ",")
	if len(specs) > 0 {
		var execErr []string
		wg := sync.WaitGroup{}
		wg.Add(len(specs))
		for _, v := range specs {
			go func(s string, c string, b string, execErr []string) {
				f := filepath.Base(s)
				args := []string{
					"run",
					"--browser",
					b,
					"--headless",
					"--spec",
					s,
					"--reporter",
					"mochawesome",
					"--reporter-options",
					fmt.Sprintf("reportFilename=%s", f),
				}
				log.Debug().Msgf("Running cypress command %s %s", "/usr/bin/cypress", strings.Join(args, " "))
				output, err := exec.Command(
					"/usr/bin/cypress",
					"run",
					"--browser",
					b,
					"--headless",
					"--spec",
					s,
					"--reporter",
					"mochawesome",
					"--reporter-options",
					fmt.Sprintf("reportFilename=%s", f),
				).Output()
				log.Debug().Msgf("Execution output %s", string(output))
				if err != nil {
					execErr = append(execErr, err.Error())
				}
				wg.Done()
			}(v, c.ConfigFile, c.Browser, execErr)
		}
		wg.Wait()
		log.Info().Msg("Program execution successful")
	}
}
