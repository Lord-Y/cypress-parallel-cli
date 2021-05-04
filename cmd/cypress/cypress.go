// Package cypress assemble all commands required to run cypress unit testing
package cypress

import (
	"fmt"
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

// Run will run cypress command
func (c *Cypress) Run() {
	fmt.Printf("run\n")
}
