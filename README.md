# cypress-parallel-cli [![CircleCI](https://circleci.com/gh/Lord-Y/cypress-parallel-cli/tree/main.svg?style=svg)](https://circleci.com/gh/Lord-Y/cypress-parallel-cli?branch=main)

`cypress-parallel-cli` is the tool that will be used by cypress-parallel to run cypress unit testing in parallel.

At the execution of each spec, it will send back the result to the API.

## Debug

To debug your code directly in the docker container, do this:

```bash
sudo docker run --rm --entrypoint="" -v $PWD:/tmp/cypress-parallel-cli -exec -ti ghcr.io/lord-y/cypress-parallel-docker-images/cypress-parallel-docker-images:9.7.0-0.1.1 bash
# install golang
cd /tmp/cypress-parallel-cli
scripts/golang.sh 1.19.2
source ~/.bashrc
# exemple of command
export CYPRESS_PARALLEL_CLI_LOG_LEVEL=debug
export CYPRESS_PARALLEL_CLI_LOG_LEVEL_WITH_CALLER=true
go run main.go cypress --repository https://github.com/cypress-io/cypress-example-kitchensink.git --branch refs/tags/v1.15.3 --specs cypress/integration/2-advanced-examples/connectors.spec.js --uid uuid --report-back
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
golangci-lint run
```

## Tests
```bash
go test -v ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
go tool cover -html=coverage.out
```

## CI

To know which `cypress`, `chrome` and `firefox` version to use for version `9.7.0`, do:
```bash
sudo docker run --rm --entrypoint="" -exec -ti cypress/included:9.7.0 bash
cypress --version
firefox --version
google-chrome --version
```
