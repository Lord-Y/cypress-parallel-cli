// Package cypress assemble all commands required to run cypress unit testing
package cypress

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"os/signal"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRun_fail(t *testing.T) {
	var c Cypress
	c.Repository = "https://github.com/cypress-io/cypress-example-kitchensink.git"
	c.Specs = "cypress/integration/examples/actions.spec.js"
	c.UniqID = "uid"
	c.Branch = "test"
	c.Username = "test"

	if os.Getenv("FATAL") == "1" {
		c.Run()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestRun_fail")
	cmd.Env = append(os.Environ(), "FATAL=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %s, want exit status 1", err.Error())
}

func TestRun(t *testing.T) {
	var c Cypress
	c.Repository = "https://github.com/cypress-io/cypress-example-kitchensink.git"
	c.Specs = "cypress/integration/examples/actions.spec.js"
	c.UniqID = "uid"

	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal(err)
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)

	go func() {
		<-sigc
		c.Run()
		signal.Stop(sigc)
	}()

	proc.Signal(os.Interrupt)
	time.Sleep(1 * time.Second)
}

func TestRun_fail_config_file(t *testing.T) {
	var (
		c      Cypress
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	c.Repository = "https://github.com/cypress-io/cypress-example-kitchensink.git"
	c.Specs = "cypress/integration/examples/actions.spec.js"
	c.UniqID = "uid"
	c.ConfigFile = "cypress_fail.json"

	if os.Getenv("FATAL") == "1" {
		c.Run()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestRun_fail_config_file")
	cmd.Env = append(os.Environ(), "FATAL=1")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %s, want exit status 1", err.Error())
}

func TestRun_success(t *testing.T) {
	assert := assert.New(t)

	var (
		c      Cypress
		stdout bytes.Buffer
	)
	c.Repository = "https://github.com/cypress-io/cypress-example-kitchensink.git"
	c.Specs = "cypress/integration/examples/actions.spec.js"
	c.UniqID = "uid"

	if os.Getenv("FATAL") == "0" {
		c.Run()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestRun_success")
	cmd.Env = append(os.Environ(), "FATAL=0")
	cmd.Stdout = &stdout
	err := cmd.Run()
	assert.NoError(err)
}

func TestRun_success_reportback(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello")
	}))
	defer ts.Close()

	var (
		c Cypress
	)
	c.Repository = "https://github.com/cypress-io/cypress-example-kitchensink.git"
	c.Specs = "cypress/integration/examples/actions.spec.js"
	c.UniqID = "uid"
	c.ReportBack = true
	c.ApiURL = ts.URL
	c.Run()
}

func TestRun_fail_reportback(t *testing.T) {
	assert := assert.New(t)
	os.Setenv("HTTP_RETRY_MAX", "1")
	defer os.Unsetenv("HTTP_RETRY_MAX")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Internal Server Error")
	}))
	defer ts.Close()

	var (
		c      Cypress
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	c.Repository = "https://github.com/cypress-io/cypress-example-kitchensink.git"
	c.Specs = "cypress/integration/examples/actions.spec.js"
	c.UniqID = "uid"
	c.ReportBack = true
	c.ApiURL = ts.URL

	if os.Getenv("FATAL") == "1" {
		c.Run()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestRun_fail_reportback")
	cmd.Env = append(os.Environ(), "FATAL=1")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	assert.NoError(err)
}

func TestReportBack(t *testing.T) {
	var (
		c Cypress
	)
	c.Repository = "https://github.com/cypress-io/cypress-example-kitchensink.git"
	c.Specs = "cypress/integration/examples/actions.spec.js"
	c.UniqID = "uid"
	c.ReportBack = true

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello")
	}))
	defer ts.Close()

	c.ApiURL = ts.URL
	c.reportBack(fmt.Errorf("Execution failed"))
}
