# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.3.0](https://github.com/Lord-Y/cypress-parallel-cli/releases/tag/v0.3.0) - 2022-10-14

### Changed
- Update requirements for cypress version >=10.0.0

## [v0.2.0](https://github.com/Lord-Y/cypress-parallel-cli/releases/tag/v0.2.0) - 2022-10-14

### Changed
- Replace deprecated func io/ioutil by io
- Upgrade to golang 1.19
- Fix git clone issues
- Enforce mochawesome reportFilename
- Improve report back to api
- Add CYPRESS_PARALLEL_CLI_LOG_LEVEL_WITH_CALLER option for logger
- Remove headless option when cypress version >= 10.0.0
- Update doc
- Uninstall all npm packages matching cypress

## [v0.1.1](https://github.com/Lord-Y/cypress-parallel-cli/releases/tag/v0.1.0) - 2021-12-29

### Changed
- Update documentation

## [v0.1.0](https://github.com/Lord-Y/cypress-parallel-cli/releases/tag/v0.1.0) - 2021-12-20

### Changed
- Update cypress-parallel api url
- Upgrade golang version

## [v0.0.5](https://github.com/Lord-Y/cypress-parallel-cli/releases/tag/v0.0.5) - 2021-06-04

### Changed
- Fix long running execution with kubernetes
- Fix cypress screen issues with Xvfb
- Drop Windows support

## [v0.0.4](https://github.com/Lord-Y/cypress-parallel-cli/releases/tag/v0.0.4) - 2021-06-02

### Changed
- Fix long running execution with kubernetes

## [v0.0.3](https://github.com/Lord-Y/cypress-parallel-cli/releases/tag/v0.0.3) - 2021-05-31

### Changed
- Fix panic error when reporting back execution status

## [v0.0.2](https://github.com/Lord-Y/cypress-parallel-cli/releases/tag/v0.0.2) - 2021-05-18

### Added
- Add context timeout while executing commands.

## [v0.0.1](https://github.com/Lord-Y/cypress-parallel-cli/releases/tag/v0.0.1) - 2021-05-15

Initial version

### Added
- Add cypress options to clone and launch cypress unit testing.
- Add version-details to print full version of the cli

## [v0.0.1-beta1](https://github.com/Lord-Y/cypress-parallel-cli/releases/tag/v0.0.1-beta1) - 2021-05-15

Initial beta version

### Added
- Add cypress options to clone and launch cypress unit testing.
- Add version-details to print full version of the cli
