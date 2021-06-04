# cypress-parallel-cli [![CircleCI](https://circleci.com/gh/Lord-Y/cypress-parallel-cli/tree/main.svg?style=svg)](https://circleci.com/gh/Lord-Y/cypress-parallel-cli?branch=main)

`cypress-parallel-cli` is the tool that will be used by cypress-parallel-api to run cypress unit testing in parallel.

At the execution of each spec, it will send back the result to the API.

## Debug

To debug your code directly in the docker container, do this:

```bash
# install golang
scripts/golang.sh 1.16.5
sudo docker run --rm --entrypoint="" -v ~/cypress-parallel-cli:/tmp/cypress-parallel-cli -exec -ti docker.pkg.github.com/lord-y/cypress-parallel-docker-images/cypress-parallel-docker-images:7.2.0-0.0.4 bash
```