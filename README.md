# cypress-parallel-cli [![CircleCI](https://circleci.com/gh/Lord-Y/cypress-parallel-cli/tree/main.svg?style=svg)](https://circleci.com/gh/Lord-Y/cypress-parallel-cli?branch=main)

`cypress-parallel-cli` is the tool that will be used by cypress-parallel to run cypress unit testing in parallel.

At the execution of each spec, it will send back the result to the API.

## Debug

To debug your code directly in the docker container, do this:

```bash
sudo docker run --rm --entrypoint="" -v $PWD:/tmp/cypress-parallel-cli -exec -ti ghcr.io/lord-y/cypress-parallel-docker-images/cypress-parallel-docker-images:7.4.0-0.1.0 bash
# install golang
scripts/golang.sh 1.17.5
source ~/.bashrc
# exemple of command
export CYPRESS_PARALLEL_CLI_LOG_LEVEL=debug
go run main.go cypress --repository https://github.com/cypress-io/cypress-example-kitchensink.git --specs cypress/integration/2-advanced-examples/connectors.spec.js --uid uuid --report-back
```

## Git hooks

Add githook like so:

```bash
git config core.hooksPath .githooks
```

## Linter
```bash
# https://golangci-lint.run/usage/install/
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```