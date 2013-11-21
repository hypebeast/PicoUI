package main

/*
 * Picoui-Chief
 *
 * Part of the PicoUi project.
 *
 * The Picoui-chief monitores and controls PicoUi applications. Moreover, it
 * provides system information about the Raspberry Pi (or any other linux
 * system) on which it is running.
 *
 * Created: 2013.11.21
 * Author: Sebastian Ruml, sebastian.ruml@gmail.com
 */

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
)

const (
	CONFIG_FILE = "picoui-chief.json"
)

var (
	config      *Config
	app_running bool = false
	running_app *exec.Cmd
)

func readProc(file string) string {
	contents, err := ioutil.ReadFile(fmt.Sprintf("/proc/%s", file))
	if err != nil {
		return ""
	}
	return string(contents[:])
}

func executeCmd(command string, args string) string {
	out, err := exec.Command(command, args).Output()
	if err != nil {
		return fmt.Sprintf("error: %s", err.Error())
	}
	return string(out[:])
}

func findApps(folder string) []AppInfo {
	var apps []AppInfo

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		fmt.Printf("no such file or directory: %s", folder)
		return apps
	}

	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			appName := file.Name()
			appPath := path.Join(folder, appName)
			appFile := path.Join(appPath, appName)

			if _, err := os.Stat(appFile); err == nil {
				apps = append(apps, AppInfo{Name: appName, Path: appPath})
			}
		}
	}

	return apps
}

func startAppHandler(w http.ResponseWriter, r *http.Request) {
	appName := r.URL.RawQuery

	if app_running {
		running_app.Process.Kill()
	}

	// Check if a app with the given name exists
	apps := findApps(config.AppsFolder)
	var foundApp AppInfo
	found := false
	for _, app := range apps {
		if app.Name == appName {
			foundApp = app
			found = true
			break
		}
	}

	if !found {
		fmt.Fprintf(w, "%s", "error: app not found")
		return
	}

	running_app = exec.Command(path.Join(foundApp.Path, foundApp.Name))
	err := running_app.Start()
	if err != nil {
		fmt.Fprintf(w, "%s", "error: can't start app")
		return
	}

	app_running = false

	fmt.Fprintf(w, "%s:%s", "ok starting app: ", foundApp.Name)
}

func killAppHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func listAppsHandler(w http.ResponseWriter, r *http.Request) {
	apps := findApps(config.AppsFolder)
	enc := json.NewEncoder(w)
	err := enc.Encode(&apps)
	if err != nil {
		fmt.Println("encoding error", err)
	}
}

func uptimeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", executeCmd("uptime", ""))
}

func ifconfigHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", executeCmd("ifconfig", ""))
}

func wHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", executeCmd("w", ""))
}

func psHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", executeCmd("ps", "-ef"))
}

func lsusbHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", executeCmd("lsusb", ""))
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", readProc("version"))
}

func meminfoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", readProc("mem"))
}

func main() {
	configFile := flag.String("config", "/etc/picoui-chief.json", "Config file")
	version := flag.Bool("version", false, "Output version and exit")

	flag.Parse()

	if *version {
		fmt.Printf("picoui-chief v%s\n", VERSION)
		fmt.Println("picoui-chief is part of the PicoUi project")
		fmt.Println("2013, Sebastian Ruml <sebastian.ruml@gmail.com>")
		os.Exit(0)
	}

	// Load the configuration
	config = loadConfig(*configFile)

	fmt.Printf("Starting picoui-chief v%s\n", VERSION)

	// Set up all handlers
	http.HandleFunc("/apps/start", startAppHandler)
	http.HandleFunc("/apps/kill", killAppHandler)
	http.HandleFunc("/apps/list", listAppsHandler)
	http.HandleFunc("/system/uptime", uptimeHandler)
	http.HandleFunc("/system/ifconfig", ifconfigHandler)
	http.HandleFunc("/system/w", wHandler)
	http.HandleFunc("/system/ps", psHandler)
	http.HandleFunc("/system/lsusb", lsusbHandler)
	http.HandleFunc("/system/version", versionHandler)
	http.HandleFunc("/system/mem", meminfoHandler)

	// Start server and listen for incoming requests
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}

func loadConfig(filename string) (c *Config) {
	if filename == "" {
		fmt.Println("Error: You should specify a config file")
		flag.Usage()
		os.Exit(-1)
	}

	text, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Can't read config file: %s\n", filename)
		os.Exit(-1)
	}

	config := &Config{}
	err = json.Unmarshal(text, &config)
	if err != nil {
		fmt.Printf("Can't parse config file: %s\n", filename)
		os.Exit(-1)
	}

	return config
}
