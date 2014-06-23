package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

var cmdNew = &Command{
	UsageLine: "new [appname]",
	Short:     "create a skeleton PicoUi application",
	Long: `
New creates a new PicoUi application with all required files.

It puts all of the files in a folder named [appname] in the current directory.

The [appname] folder has the following contents:

	|- main.go
	|- static
	  |- index.html
	  |- js
	  |- css
	  |- fonts

For example:

	picoui new helloworld
`,
}

func init() {
	cmdNew.Run = newApp
}

func newApp(cmd *Command, args []string) {
	if len(args) == 0 {
		errorf("No app name given.\nRun 'picoui help new' for usage.\n")
	}

	curpath, _ := os.Getwd()
	errorIfGopathIsNotSet()

	appName := args[0]

	// Absolute paths are not allowed
	if path.IsAbs(appName) {
		errorf("Abort: '%s' looks like an absolute path. Please provide an application name instead.", appName)
	}

	gopath := os.Getenv("GOPATH")

	// Check if current path is in GOPATH
	if !isPathInGopath(curpath) {
		errorf("Abort: Unable to create an application outside of GOPATH '%s'\n"+
			"Change your work directory by 'cd %s%ssrc'\n", gopath, gopath, filepath.Separator)
	}

	gosrcpath := path.Join(gopath, "src")

	// Check if the picoui-lib source can be found
	srcpath := filepath.Join(gosrcpath, PICOUI_IMPORT_PATH)
	if _, err := os.Stat(srcpath); os.IsNotExist(err) {
		errorf("Abort: Could not find PicoUi source code: %s\n", err)
	}

	appPath := path.Join(curpath, appName)

	// Check if the app folder already exists
	if _, err := os.Stat(appPath); !os.IsNotExist(err) {
		errorf("Abort: '%s' already exists.\n", appPath)
	}

	fmt.Fprintf(os.Stdout, "Creating application...\n\n")

	appInfo := AppInfo{Name: appName}

	os.MkdirAll(appPath, 0755)
	copyDir(appPath, filepath.Join(srcpath, "skeleton"), appInfo)

	fmt.Fprintf(os.Stdout, "Your application is ready:\n")
	fmt.Fprintf(os.Stdout, "\t%s\n\n", appPath)
	fmt.Fprintf(os.Stdout, "You can start it with:\n")
	fmt.Fprintf(os.Stdout, "\tcd %s\n\tpicoui run\n", appName)
}
