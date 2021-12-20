// Package cypress assemble all commands required to run cypress unit testing
package cypress

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Lord-Y/cypress-parallel-cli/git"
	"github.com/Lord-Y/cypress-parallel-cli/httprequests"
	"github.com/Lord-Y/cypress-parallel-cli/logger"
	"github.com/rs/zerolog/log"
)

var apiURI = "/api/v1/executions/update"

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
	Timeout    int    // Timeout after which the program will exit with error
}

func init() {
	logger.SetLoggerLogLevel()
}

// Run will run cypress command
func (c *Cypress) Run() {
	var (
		gc git.Repository
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Timeout)*time.Minute)
	defer cancel()

	gc.Repository = c.Repository
	gc.Username = c.Username
	gc.Password = c.Password
	gc.Branch = c.Branch

	gitdir, err := gc.Clone()
	if err != nil {
		c.reportBack(err, "")
		log.Fatal().Err(err).Msg("Error occured while cloning git repository")
		return
	}
	defer os.RemoveAll(gitdir)
	log.Debug().Msgf("Git temp dir %s", gitdir)

	if ctx.Err() == context.DeadlineExceeded {
		c.reportBack(ctx.Err(), "")
		log.Fatal().Err(ctx.Err()).Msgf("Execution timeout reached after %d minute(s)", c.Timeout)
		return
	}

	err = os.Chdir(gitdir)
	if err != nil {
		c.reportBack(err, "")
		log.Fatal().Err(err).Msg("Error occured while chdir git repository")
		return
	}

	if c.ConfigFile != "" {
		var info os.FileInfo
		if info, err = os.Stat(fmt.Sprintf("%s/%s", gitdir, c.ConfigFile)); os.IsNotExist(err) {
			c.reportBack(err, "")
			log.Fatal().Err(err).Msgf("Error occured while checking config file %s", c.ConfigFile)
			return
		}
		if !info.Mode().IsRegular() {
			c.reportBack(err, "")
			log.Fatal().Err(err).Msgf("Error occured while checking config file %s", c.ConfigFile)
			return
		}
	}

	execUninstallCmd := exec.CommandContext(
		ctx,
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
		c.reportBack(err, "")
		log.Fatal().Err(err).Msg("Error occured while forcing uninstall of local cypress package")
		return
	}

	output, err := exec.CommandContext(
		ctx,
		"npm",
		"install",
	).Output()
	log.Debug().Msgf("NPM install output %s", string(output))

	if err != nil {
		c.reportBack(fmt.Errorf("%s | %s", string(output), err), "")
		log.Fatal().Err(fmt.Errorf("%s | %s", string(output), err)).Msgf("Error occured while installing user packages")
		return
	}

	output, err = exec.CommandContext(
		ctx,
		"npm",
		"install",
		"--save-dev",
		"mochawesome",
	).Output()
	log.Debug().Msgf("Mochawesome install output %s", string(output))

	if err != nil {
		c.reportBack(fmt.Errorf("%s | %s", string(output), err), "")
		log.Fatal().Err(fmt.Errorf("%s | %s", string(output), err)).Msgf("Error occured while installing mochawesome")
		return
	}

	specs := strings.Split(c.Specs, ",")
	wg := sync.WaitGroup{}
	for i, spec := range specs {
		wg.Add(1)
		go func(i int, spec string) {
			defer wg.Done()
			// https://docs.cypress.io/guides/continuous-integration/introduction#Xvfb
			screen := fmt.Sprintf(":%d", 99+i)
			cmd := exec.Command("sh", "-c", fmt.Sprintf("Xvfb %s &", screen))
			log.Debug().Msgf("Execution output of screen command %s", cmd.String())
			err := cmd.Run()
			if err != nil {
				log.Error().Err(err).Msgf("Fail to execute Xvfb command %s", err.Error())
				c.reportBack(err, spec)
				return
			}

			f := filepath.Base(spec)
			args := []string{
				"run",
				"--browser",
				c.Browser,
				"--headless",
				"--spec",
				spec,
				"--reporter",
				"mochawesome",
				"--reporter-options",
				fmt.Sprintf("reportFilename=%s", f),
			}
			log.Debug().Msgf("Running cypress command %s %s", "cypress", strings.Join(args, " "))

			process := exec.CommandContext(
				ctx,
				"cypress",
				"run",
				"--browser",
				c.Browser,
				"--headless",
				"--spec",
				spec,
				"--reporter",
				"mochawesome",
				"--reporter-options",
				fmt.Sprintf("reportFilename=%s", f),
			)
			process.Env = append(os.Environ(), fmt.Sprintf("DISPLAY=%s", screen))
			process.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
			go func() {
				<-ctx.Done()
				if ctx.Err() == context.DeadlineExceeded {
					syscall.Kill(-process.Process.Pid, syscall.SIGKILL)
					c.reportBack(ctx.Err(), spec)
					log.Error().Err(ctx.Err()).Msgf("Execution timeout reached after %d minute(s)", c.Timeout)
					return
				}
			}()
			output, err := process.Output()
			log.Debug().Msgf("Execution output %s", string(output))
			if err != nil {
				log.Error().Err(err).Msgf("Fail to execute cypress command %s", string(output))
				c.reportBack(fmt.Errorf("%s | %s", string(output), err), spec)
				return
			}
			log.Debug().Msgf("Reporting back result to %s", fmt.Sprintf("%s%s", c.ApiURL, apiURI))
			result := fmt.Sprintf("%s/mochawesome-report/%s.json", gitdir, strings.TrimSuffix(f, ".js"))
			of, err := os.Open(result)
			if err != nil {
				log.Error().Err(err).Msgf("Fail to open file %s", result)
				c.reportBack(err, spec)
				return
			}
			defer of.Close()
			fo, err := io.ReadAll(of)
			if err != nil {
				log.Error().Err(err).Msgf("Fail to read file %s content", result)
				c.reportBack(err, spec)
				return
			}
			headers := make(map[string]string)
			headers["Content-Type"] = "application/x-www-form-urlencoded"
			buf := new(bytes.Buffer)
			if err := json.Compact(buf, fo); err != nil {
				log.Error().Err(err).Msg("Fail to compact json result")
				c.reportBack(err, spec)
				return
			}

			payload := fmt.Sprintf("result=%s", hex.EncodeToString(buf.Bytes()))
			payload += "&executionStatus=DONE"
			payload += fmt.Sprintf("&uniqId=%s", c.UniqID)
			payload += fmt.Sprintf("&branch=%s", c.Branch)
			payload += fmt.Sprintf("&spec=%s", spec)
			payload += fmt.Sprintf("&encoded=%s", "true")

			_, _, err = httprequests.PerformRequests(headers, "POST", fmt.Sprintf("%s%s", c.ApiURL, apiURI), payload, "")
			if err != nil {
				log.Error().Err(err).Msg("Fail to report back result")
			}
		}(i, spec)
	}
	wg.Wait()
	log.Info().Msg("Program execution successful")
}

func (c *Cypress) reportBack(err error, spec string) {
	if c.ReportBack {
		headers := make(map[string]string)
		headers["Content-Type"] = "application/x-www-form-urlencoded"
		if spec != "" {
			payload := url.Values{}
			payload.Set("result", "{}")
			payload.Set("executionStatus", "FAILED")
			payload.Set("uniqId", c.UniqID)
			payload.Set("branch", c.Branch)
			payload.Set("spec", spec)
			payload.Set("executionErrorOutput", err.Error())

			_, _, err = httprequests.PerformRequests(headers, "POST", fmt.Sprintf("%s%s", c.ApiURL, apiURI), payload.Encode(), "")
			if err != nil {
				log.Error().Err(err).Msg("Fail to report back result")
				return
			}
		} else {
			specs := strings.Split(c.Specs, ",")
			for _, spec := range specs {
				payload := url.Values{}
				payload.Set("result", "{}")
				payload.Set("executionStatus", "FAILED")
				payload.Set("uniqId", c.UniqID)
				payload.Set("branch", c.Branch)
				payload.Set("spec", spec)
				payload.Set("executionErrorOutput", err.Error())

				_, _, err = httprequests.PerformRequests(headers, "POST", fmt.Sprintf("%s%s", c.ApiURL, apiURI), payload.Encode(), "")
				if err != nil {
					log.Error().Err(err).Msg("Fail to report back result")
					return
				}
			}
		}
	}
}
