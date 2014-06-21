package main

import (
	"os"
)

var cmdPublish = &Command{
	UsageLine: "publish [-h host] [-p path] [-u user] [-p password] [-n]",
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

TODO: -p flag

TODO: -n flag (By default the command cross-compiles for the Raspberry Pi.)

The following environment variables can used instead of the flags:

  - PICO_USER: The username that is used for login on the remote host
  - PICO_PASS: The password
  - PICO_HOST:
  - PICO_PATH:

For example:

	picoui publish -h raspberry.local -p /opt/picoui/apps -u rasp -p berry
`,
}

func init() {
	cmdPublish.Run = publishApp

	// TODO: Add flags
}

func publishApp(cmd *Command, args []string) {
	if len(args) < 2 {
		errorf("Abort: No host and remote path given.\nRun 'picoui help publish' for more information.")
	}

	if len(args) > 4 {
		errorf("Abort: Do many arguments.\nRun 'picoui help publish' for more information.")
	}

	curpath, _ := os.Getwd()

	Debugf("current path: %s", curpath)

	// Check if GOPATH is set
	errorIfGopathIsNotSet()

	// Check if current path is in GOPATH
	if !isPathInGopath(curpath) {
		errorf("Abort: Unable to run application outside of GOPATH '%s'", os.Getenv("GOPATH"))
	}

	// credentialsProvided := false
	// username := ""
	// password := ""

	// // Check if there is a username and password specified
	// if len(args) > 2 {
	// 	username = args[2]
	// 	password = args[3]
	// 	credentialsProvided = true
	// } else {
	// 	if os.Getenv("PICUI_USER") != "" && os.Getenv("PICOUI_PASSWORD") != "" {
	// 		username = os.Getenv("PICOUI_USER")
	// 		password = os.Getenv("PICOUI_PASSWORD")
	// 		credentialsProvided = true
	// 	}
	// }

	// appname := path.Base(curpath)
	// fmt.Fprintf(os.Stdout, "Uses '%s' as appname\n", appname)

	// Check for
	// TODO: Build the executable

	// TODO: Copy the executable and static files to the remote host
}
