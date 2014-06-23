package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

var cmdBuild = &Command{
	UsageLine: "build [-pi]",
	Short:     "build the application",
	Long: `
Build the PicoUi application. This command must be executed in the application
directory.

If the -pi flag is specified the application will be cross compiled for the Raspberry Pi.

Example:

    picoui build
`,
}

func init() {
	cmdBuild.Run = runBuild

	cmdBuild.Flag.BoolVar(&buildForPi, "pi", false, "")
}

var buildForPi bool // -p flag

func runBuild(cmd *Command, args []string) {
	curpath, _ := os.Getwd()
	Debugf("current path: %s", curpath)

	// Check if GOPATH is set
	errorIfGopathIsNotSet()

	// Check if current path is in GOPATH
	if !isPathInGopath(curpath) {
		errorf("Abort: Unable to run application outside of GOPATH '%s'", os.Getenv("GOPATH"))
	}

	appname := path.Base(curpath)
	Debugf("app name: %s", appname)

	fmt.Fprintf(os.Stdout, "Building application '%s'...\n", appname)

	app := exec.Command("go", "build", "-o", appname, "./...")

	if buildForPi {
		Debugf("Cross compiling for RaspberryPi")
		app.Env = append(app.Env, "GOARCH=arm")
		app.Env = append(app.Env, "GOARM=5")
		app.Env = append(app.Env, "GOOS=linux")
	}

	app.Dir = curpath
	err := app.Start()
	if err != nil {
		errorf("Abort: %s", err)
	}
	app.Wait()
	fmt.Fprintf(os.Stdout, "Finished building\n")
}
