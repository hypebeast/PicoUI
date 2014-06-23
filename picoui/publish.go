package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

var cmdPublish = &Command{
	UsageLine: "publish [-h host] [-dir path] [-u user] [-p password] [-pi]",
	Short:     "publish a PicoUi application to a remote host",
	Long: `
Publish the PicoUi application to a remote host where picoui-chief is installed.
This command uses scp to copy the application to the remote host.

The best way to use this command is to set up SSH public/private keys for your
development and remote host. See http://www.ece.uci.edu/~chou/ssh-key.html for
more information. Otherwise, you need to specify a username (-u flag) and
password (-p flag) as arguments. You can also set the PICOUI_USER and
PICOUI_PASSWORD environment variables to provide the username and password.

TODO: -h flag

TODO: -dir flag

TODO: -u flag

TODO: -p flag

If the -pi flag is specified the application will be cross compiled for the Raspberry Pi.

The following environment variables can be set:

  - PICO_USER: The username that is used for authentication
  - PICO_PASS: The password that is used for authentication
  - PICO_HOST: The remote host
  - PICO_PATH: The remote installation path (default path: /opt/picoui/apps)

For example:

	picoui publish -h raspberry.local -u rasp -p berry
`,
}

var (
	defaultRemotePath = "/opt/picoui/apps"
)

func init() {
	cmdPublish.Run = publishApp

	cmdPublish.Flag.StringVar(&hostFlag, "h", "", "")
	cmdPublish.Flag.StringVar(&remotePathFlag, "dir", "", "")
	cmdPublish.Flag.StringVar(&userFlag, "u", "", "")
	cmdPublish.Flag.StringVar(&passwordFlag, "p", "", "")
	cmdPublish.Flag.BoolVar(&buildForPi, "pi", false, "")
}

var hostFlag string       // -h flag
var remotePathFlag string // -dir flag
var userFlag string       // -u flag
var passwordFlag string   // -p flag

func publishApp(cmd *Command, args []string) {
	curpath, _ := os.Getwd()
	Debugf("current path: %s", curpath)

	// Check if GOPATH is set
	errorIfGopathIsNotSet()

	// Check if current path is in GOPATH
	if !isPathInGopath(curpath) {
		errorf("Abort: Unable to run application outside of GOPATH '%s'", os.Getenv("GOPATH"))
	}

	// Check the flags and environment variables
	host := os.Getenv("PICO_HOST")
	remotePath := os.Getenv("PICO_PATH")
	user := os.Getenv("PICO_USER")
	password := os.Getenv("PICO_PASS")

	if hostFlag != "" {
		host = hostFlag
	}

	if remotePathFlag != "" {
		remotePath = remotePathFlag
	}

	if userFlag != "" {
		user = userFlag
	}

	if passwordFlag != "" {
		password = passwordFlag
	}

	if host == "" {
		errorf("Abort: No remote host specified.")
	}

	if remotePath == "" {
		remotePath = defaultRemotePath
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
	err = app.Wait()
	if err != nil {
		errorf("Abort: %s", err)
	}

	fmt.Fprintf(os.Stdout, "Building package...\n")

	packageName := appname + ".tar.gz"
	runCommand(curpath, "tar", "cvzf", packageName, appname, "static")

	fmt.Fprintf(os.Stdout, "Publishing package to '%s:%s'...\n", host, remotePath)

	var usePubKey bool
	var sshpassAvailable bool
	// If no password is given, try to use public key authentication
	if password == "" {
		usePubKey = true
	} else {
		usePubKey = false

		// Check if sshpass is available
		app = exec.Command("sshpass", "-V")
		err = app.Start()
		if err != nil {
			errorf("Abort: %s", err)
		}
		err = app.Wait()
		if err != nil {
			// sshpass isn't available
			sshpassAvailable = false
		} else {
			sshpassAvailable = true
		}
	}

	var connectionString string
	if user == "" {
		connectionString = fmt.Sprintf("%s", host)
	} else {
		connectionString = fmt.Sprintf("%s@%s", user, host)
	}

	if !usePubKey && sshpassAvailable {
		runCommand(curpath, "sshpass", "-p", password, "scp", packageName, connectionString+":"+remotePath)
	} else {
		runCommand(curpath, "scp", packageName, connectionString+":"+remotePath)
	}

	runCommand(curpath, "rm", packageName)

	// Unpack the package on the remote host
	command := "cd " + remotePath + "; mkdir " + appname + "; tar xvzf " + packageName + " -C " + appname + "; rm " + packageName
	if !usePubKey && sshpassAvailable {
		runCommand(curpath, "sshpass", "-p", password, "ssh", connectionString, command)
	} else {
		runCommand(curpath, "ssh", connectionString, command)
	}

	fmt.Fprintf(os.Stdout, "Done\n")
}
