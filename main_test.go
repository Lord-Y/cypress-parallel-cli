package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	cli "github.com/urfave/cli/v2"
)

func TestCypress(t *testing.T) {
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
				"--specs",
				"cypress/integration/examples/actions.spec.js",
			},
			fail: true,
		},
		{
			cliArgs: []string{
				"cypress-parallel-cli",
				"cypress",
				"--repository",
				"https://github.com/cypress-io/cypress-example-kitchensink.git",
				"--specs",
				"cypress/integration/examples/actions.spec.js",
				"--uid",
				"uid",
			},
			fail: false,
		},
	}

	app := cli.NewApp()
	for _, tc := range tests {
		app.Writer = ioutil.Discard
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
