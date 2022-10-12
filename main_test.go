package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"os/signal"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	cli "github.com/urfave/cli/v2"
)

func TestMain(t *testing.T) {
	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)

	go func() {
		<-sigc
		os.Args = []string{
			"cypress",
			"-h",
		}
		main()
		signal.Stop(sigc)
	}()

	_ = proc.Signal(os.Interrupt)
	time.Sleep(1 * time.Second)
}

func TestMain_fail(t *testing.T) {
	assert := assert.New(t)
	if os.Getenv("FATAL") == "1" {
		os.Args = []string{
			"cypress",
			"bad",
		}
		main()
		return
	}
	cmd := exec.Command(os.Args[0], "cypress", "bad", "-test.run=TestMain_fail")
	cmd.Env = append(os.Environ(), "FATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	assert.Error(err)
}

func TestMainCypress(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		cliArgs []string
		fail    bool
	}{
		{
			cliArgs: []string{
				"cypress-parallel-cli",
				"cypress",
				"--repository",
				"test",
			},
			fail: true,
		},
		{
			cliArgs: []string{
				"cypress-parallel-cli",
				"cypress",
				"--repository",
				"https://github.com/cypress-io/cypress-example-kitchensink.git",
				"--branch",
				"refs/tags/v1.15.3",
				"--specs",
				"cypress/integration/2-advanced-examples/connectors.spec.js",
			},
			fail: true,
		},
		{
			cliArgs: []string{
				"cypress-parallel-cli",
				"cypress",
				"--repository",
				"https://github.com/cypress-io/cypress-example-kitchensink.git",
				"--branch",
				"refs/tags/v1.15.3",
				"--specs",
				"cypress/integration/2-advanced-examples/connectors.spec.js",
				"--uid",
				"uid",
			},
			fail: false,
		},
	}

	app := cli.NewApp()
	for _, tc := range tests {
		app.Writer = io.Discard
		app.Commands = []*cli.Command{
			cypress,
		}
		err := app.Run(tc.cliArgs)
		if tc.fail {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
	}
}

func TestMainCypressReportBack(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello")
	}))
	defer ts.Close()

	tests := []struct {
		cliArgs []string
		fail    bool
	}{
		{
			cliArgs: []string{
				"cypress-parallel-cli",
				"cypress",
				"--repository",
				"https://github.com/cypress-io/cypress-example-kitchensink.git",
				"--branch",
				"refs/tags/v1.15.3",
				"--specs",
				"cypress/integration/2-advanced-examples/connectors.spec.js",
				"--uid",
				"uid",
				"--rp",
				"--api-url",
			},
			fail: false,
		},
	}

	app := cli.NewApp()
	for _, tc := range tests {
		app.Writer = io.Discard
		app.Commands = []*cli.Command{
			cypress,
		}
		err := app.Run(append(tc.cliArgs, ts.URL))
		if tc.fail {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
	}
}

func TestVersionDetails(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		cliArgs []string
		fail    bool
	}{
		{
			cliArgs: []string{
				"cypress-parallel-cli",
				"version-details",
			},
			fail: false,
		},
	}

	app := cli.NewApp()
	for _, tc := range tests {
		app.Writer = io.Discard
		app.Commands = []*cli.Command{
			versionDetails,
		}
		err := app.Run(tc.cliArgs)
		if tc.fail {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
	}
}
