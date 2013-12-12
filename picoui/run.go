package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
)

var cmdRun = &Command{
	UsageLine: "run",
	Short:     "run a PicoUi application",
	Long: `
Run the PicoUi application. To run the application you must be in the directory
where the application resides.

For example:

	picoui run
`,
	Run: runApp,
}

func runApp(args []string) {
	if len(args) != 0 {
		errorf("Abort: Do many arguments given.\nRun 'picoui help run' for more information.")
	}

	curpath, _ := os.Getwd()

	Debugf("current path: %s", curpath)

	// Check if GOPATH is set
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		errorf("Abort: GOPATH environment variable is not set.\n" +
			"Please refer to http://golang.org/doc/code.html to configure your Go environment.")
	}

	// Check if current path is in GOPATH
	if !isPathInGopath(curpath) {
		errorf("Abort: Unable to run application outside of GOPATH '%s'", gopath)
	}

	appname := path.Base(curpath)
	fmt.Fprintf(os.Stdout, "Uses '%s' as appname\n", appname)

	// Find all go files and check if at least one go file is found
	var gofiles []string
	files, _ := ioutil.ReadDir(curpath)
	for _, file := range files {
		if !file.IsDir() {
			if strings.HasSuffix(file.Name(), ".go") {
				gofiles = append(gofiles, file.Name())
				break
			}
		}
	}

	if len(gofiles) < 1 {
		errorf("Abort: No go file found!")
	}

	fmt.Fprintf(os.Stdout, "Starting application '%s'... (Press 'Ctrl-C' to stop it)\n", appname)

	app := exec.Command("go", "run", strings.Join(gofiles, " "))
	app.Dir = curpath
	err := app.Start()
	if err != nil {
		errorf("Abort: %s", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			app.Process.Kill()
		}
	}()

	app.Wait()

	fmt.Fprintf(os.Stdout, "'%s' stopped\n", appname)
}
