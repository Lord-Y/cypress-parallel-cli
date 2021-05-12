// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/magefile/mage/mg"
)

var (
	packageName  = "github.com/Lord-Y/cypress-parallel-cli"
	app          = "cypress-parallel-cli"
	architecture = map[string]string{
		"linux":   "amd64",
		"darwin":  "amd64",
		"windows": "amd64",
	}
)

var ldflags = fmt.Sprintf(
	"-X %s/cmd.Version=%s -X %s/cmd.revision=%s -X '%s/cmd.buildDate=%s' -X '%s/cmd.goVersion=%s'",
	packageName,
	os.Getenv("BUILD_VERSION"),
	packageName,
	os.Getenv("BUILD_REVISION"),
	packageName,
	time.Now().Format("2006-01-02T15:04:05Z0700"),
	packageName,
	runtime.Version(),
)

// Build Build binaries depending on os/architectures
func Build() (err error) {
	mg.Deps(InstallDeps)

	wg := sync.WaitGroup{}
	wg.Add(len(architecture))
	for oses, arch := range architecture {
		go func(oses, arch string) (err error) {
			appName := fmt.Sprintf("%s_%s_%s", app, oses, arch)
			os.Remove(appName)
			fmt.Printf("Building %s ...\n", appName)
			cmd := exec.Command(
				"go",
				"build",
				"-ldflags",
				ldflags,
				"-o",
				appName,
				".",
			)
			fmt.Printf("command %s\n", cmd.String())
			cmd.Env = append(
				os.Environ(),
				fmt.Sprintf("GOOS=%s", oses),
				fmt.Sprintf("GOARCH=%s", arch),
			)
			err = cmd.Run()
			wg.Done()
			return
		}(oses, arch)
	}
	wg.Wait()
	return
}

// Manage your deps, or running package managers.
func InstallDeps() error {
	os.Setenv("GO111MODULE", "on")
	fmt.Println("Installing Deps...")
	cmd := exec.Command("go", "mod", "download")
	return cmd.Run()
}

// Clean Remove all generated binaries
func Clean() {
	wg := sync.WaitGroup{}
	wg.Add(len(architecture))
	for oses, arch := range architecture {
		go func(oses, arch string) {
			appName := fmt.Sprintf("%s_%s_%s", app, oses, arch)
			fmt.Printf("Cleaning %s ...\n", appName)
			os.Remove(appName)
			wg.Done()
		}(oses, arch)
	}
	wg.Wait()
}
