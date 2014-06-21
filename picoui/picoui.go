package main

/*
 * picoui
 *
 * Part of the PicoUi project.
 *
 * Command line tool for working with PicoUi applications
 *
 * Created: 2013.12.11
 * Author: Sebastian Ruml, sebastian.ruml@gmail.com
 */

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
)

// Command defines a command for PicoUi. The structure was copied from the
// 'go' command and the 'revel' project.
type Command struct {
	Run         func(cmd *Command, args []string)
	UsageLine   string
	Short       string
	Long        string
	Flag        flag.FlagSet // Flag is a set of flags specific to this command
	CustomFlags bool         // CustomFlags indicates that the command will do its own flag parsing
}

// Name returns the command name. The name is the first word in the usage line.
func (cmd *Command) Name() string {
	name := cmd.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

// Usage returns the usage line of the command.
func (cmd *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", cmd.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(cmd.Long))
	os.Exit(2)
}

// Available commands
var commands = []*Command{
	cmdNew,
	cmdRun,
	cmdPublish,
	cmdBuild,
}

func main() {
	fmt.Fprintf(os.Stdout, header)
	flag.Usage = func() { usage(2) }
	flag.Parse()
	args := flag.Args()

	// If no command is specified, print the usage text
	if len(args) < 1 {
		usage(2)
	}

	// Print the help text for a command
	if args[0] == "help" {
		help(args[1:])
		return
	}

	// Commands use panic to abort execution when something goes wrong.
	// Panics are logged at the point of error. Ignore those.
	defer func() {
		if err := recover(); err != nil {
			if _, ok := err.(LoggedError); !ok {
				// This panic was not expected / logged.
				panic(err)
			}
			os.Exit(1)
		}
	}()

	arg := args[0]

	for _, cmd := range commands {
		if cmd.Name() == arg && cmd.Run != nil {
			cmd.Flag.Usage = func() { cmd.Usage() }
			if cmd.CustomFlags {
				args = args[1:]
			} else {
				cmd.Flag.Parse(args[1:])
				args = cmd.Flag.Args()
			}
			cmd.Run(cmd, args)
			return
		}
	}

	errorf("unknown command %q\nRun 'revel help' for usage.\n", args[0])
}

const header = `#
# picoui --> http://github.com/hypebeast/picoui
#
`

var usageTemplate = `usage: picoui command [arguments]

The commands are:
{{range .}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}

Use "picoui help [command]" for more information.
`

var helpTemplate = `usage: picoui {{.UsageLine}}
{{.Long}}`

func errorf(format string, args ...interface{}) {
	// Ensure the user's command prompt starts on the next line
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, args...)
	panic(LoggedError{})
}

// usage prints the usage text.
func usage(exitCode int) {
	tmpl(os.Stderr, usageTemplate, commands)
	os.Exit(exitCode)
}

// help prints the usage text for a command.
func help(args []string) {
	// print the usage text and exit if no command is specified
	if len(args) < 1 {
		usage(2)
	}

	if len(args) != 1 {
		fmt.Fprintf(os.Stdout, "usage: picoui help command\n\nToo many arguments given.\n")
		os.Exit(2)
	}

	arg := args[0]

	for _, cmd := range commands {
		if cmd.Name() == arg {
			tmpl(os.Stdout, helpTemplate, cmd)
			return
		}
	}

	errorf("unknown command %q\nRun 'revel help' for usage.\n", args[0])
}

// tmpl renders a template with the given writer, template and the data.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}
