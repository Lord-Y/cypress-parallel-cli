// Package cmd manage all commands required to launch cypress-parallel-cli
package cmd

import (
	"github.com/Lord-Y/cypress-parallel-cli/cmd/cypress"
)

// list of vars that will be use by cli
var (
	// cypress struct
	cmd       cypress.Cypress
	revision  string
	buildDate string
	goVersion string
)

var Version = "0.0.1"
