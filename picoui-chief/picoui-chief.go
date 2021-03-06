package main

/*
 * Picoui-Chief
 *
 * Part of the PicoUi project.
 *
 * The Picoui-chief monitores and controls PicoUi applications. Moreover, it
 * provides some basic system information about the Raspberry Pi (or any other linux
 * system) on which it is running.
 *
 * Created: 2013.11.21
 * Author: Sebastian Ruml, sebastian.ruml@gmail.com
 */

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	CONFIG_FILE = "picoui-chief.json"
)

var (
	config         *Config
	configFileName string
	app_running    bool = false
	app_info       AppInfo
	running_app    *exec.Cmd
	prevIdle       uint64
	prevTotal      uint64
	cpuUtilization uint64
)

func readProc(file string) string {
	contents, err := ioutil.ReadFile(fmt.Sprintf("/proc/%s", file))
	if err != nil {
		return "error"
	}
	return string(contents[:])
}

func executeCmd(command string, args ...string) string {
	out, err := exec.Command(command, args...).Output()
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

func startApplication(appName string) error {
	if len(appName) < 1 {
		return errors.New("no app name")
	}

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
		return errors.New("app not found")
	}

	running_app = exec.Command(path.Join(foundApp.Path, foundApp.Name))
	running_app.Dir = foundApp.Path
	err := running_app.Start()
	if err != nil {
		return errors.New("can't start app")
	}

	app_running = true
	app_info.Name = foundApp.Name
	app_info.Path = foundApp.Path

	return nil
}

func getCpu() []uint64 {
	output := readProc("stat")
	reader := bufio.NewReader(bytes.NewBuffer([]byte(output)))
	var cpuStats []uint64
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		if string(line[0:4]) == "cpu " {
			fields := strings.Fields(string(line))

			var intVal uint64
			for i := 1; i <= 8; i++ {
				intVal, _ = strconv.ParseUint(fields[i], 10, 64)
				cpuStats = append(cpuStats, intVal)
			}

			break
		}
	}

	return cpuStats
}

func calculateCpuUsage() {
	stats := getCpu()

	var total uint64 = 0
	for _, value := range stats {
		total += value
	}

	idle := stats[3]
	diff_idle := idle - prevIdle
	diff_total := total - prevTotal
	cpuUtilization = (1000*(diff_total-diff_idle)/diff_total + 5) / 10

	prevIdle = idle
	prevTotal = total
}

func startAppHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "%s", "error: method not supported")
		return
	}

	vals := r.URL.Query()
	appName := vals["name"][0]

	// Check for autostart parameter
	if len(vals["autostart"]) > 0 {
		// Set the autostart value for this app
		config.Autostart = appName
		saveConfig(config, configFileName)
	}

	err := startApplication(appName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", "error starting app")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s: %s", "ok starting app", appName)
}

func killAppHandler(w http.ResponseWriter, r *http.Request) {
	if app_running {
		running_app.Process.Kill()
		app_running = false

		fmt.Fprintf(w, "%s", "ok killing app")
		return
	}

	fmt.Fprintf(w, "error: no running app found")
}

func setAutostartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "%s", "error: method not supported")
		return
	}

	vals := r.URL.Query()

	if len(vals["appname"]) < 1 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", "error")
		return
	}

	config.Autostart = vals["appname"][0]
	saveConfig(config, configFileName)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", "ok")
}

func disableAutostartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "%s", "error: method not supported")
		return
	}

	config.Autostart = ""
	saveConfig(config, configFileName)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", "ok")
}

func listAppsHandler(w http.ResponseWriter, r *http.Request) {
	apps := findApps(config.AppsFolder)
	enc := json.NewEncoder(w)
	err := enc.Encode(&apps)
	if err != nil {
		fmt.Println("encoding error", err)
	}
}

func appStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	enc := json.NewEncoder(w)
	if app_running {
		var info map[string]interface{}
		info["status"] = "running"
		info["info"] = app_info
		err := enc.Encode(&info)
		if err != nil {
			fmt.Println("encoding error", err)
		}
	} else {
		var info map[string]interface{}
		info["status"] = "stopped"
		err := enc.Encode(&info)
		if err != nil {
			fmt.Println("encoding error", err)
		}
	}
}

func uptimeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", executeCmd("uptime", ""))
}

func procUptimeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", readProc("uptime"))
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
	fmt.Fprintf(w, "%s", readProc("meminfo"))
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "pong")
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	err := exec.Command("sudo", "shutdown", "-h", "now").Start()
	if err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
	} else {
		fmt.Fprintf(w, "%s", "ok")
	}
}

func rebootHandler(w http.ResponseWriter, r *http.Request) {
	err := exec.Command("sudo", "reboot").Start()
	if err != nil {
		fmt.Fprintf(w, "error: %s", err.Error())
	}
	fmt.Fprintf(w, "%s", "ok")
}

func loadHandler(w http.ResponseWriter, r *http.Request) {
	output := executeCmd("uptime", "")
	parts := strings.Split(output, " ")
	load := strings.Join(parts[len(parts)-3:len(parts)], " ")

	fmt.Fprintf(w, "%s", strings.Replace(load, ",", "", -1))
}

func statHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", readProc("stat"))
}

func cpuHandler(w http.ResponseWriter, r *http.Request) {
	stats := getCpu()
	fmt.Fprintf(w, "%d %d %d %d %d %d %d %d", stats[0], stats[1], stats[2], stats[3], stats[4], stats[5], stats[6], stats[7])
}

func cpuUtilHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%d", cpuUtilization)
}

func main() {
	configFile := flag.String("config", "/etc/picoui-chief.json", "Config file")
	version := flag.Bool("version", false, "Output version and exit")

	flag.Parse()

	if *version {
		fmt.Printf("picoui-chief v%s\n", VERSION)
		fmt.Println("picoui-chief is part of the PicoUi project")
		fmt.Println("2013-2014, Sebastian Ruml <sebastian.ruml@gmail.com>")
		os.Exit(0)
	}

	// Load the configuration
	configFileName = *configFile
	config = loadConfig(configFileName)

	fmt.Printf("Starting picoui-chief v%s\n", VERSION)

	// Set up all handlers
	http.HandleFunc("/apps", listAppsHandler)
	http.HandleFunc("/apps/status", appStatusHandler)
	http.HandleFunc("/apps/start", startAppHandler)
	http.HandleFunc("/apps/kill", killAppHandler)
	http.HandleFunc("/apps/setAutostart", setAutostartHandler)
	http.HandleFunc("/apps/disableAutostart", disableAutostartHandler)
	http.HandleFunc("/system/uptime", uptimeHandler)
	http.HandleFunc("/system/proc_uptime", procUptimeHandler)
	http.HandleFunc("/system/ifconfig", ifconfigHandler)
	http.HandleFunc("/system/w", wHandler)
	http.HandleFunc("/system/ps", psHandler)
	http.HandleFunc("/system/lsusb", lsusbHandler)
	http.HandleFunc("/system/version", versionHandler)
	http.HandleFunc("/system/mem", meminfoHandler)
	http.HandleFunc("/system/shutdown", shutdownHandler)
	http.HandleFunc("/system/reboot", rebootHandler)
	http.HandleFunc("/system/load", loadHandler)
	http.HandleFunc("/system/stat", statHandler)
	http.HandleFunc("/system/cpu", cpuHandler)
	http.HandleFunc("/system/cpu_usage", cpuUtilHandler)
	http.HandleFunc("/ping", pingHandler)

	// Start the CPU usage calculation
	ticker := time.NewTicker(time.Millisecond * 1000)
	go func() {
		for _ = range ticker.C {
			calculateCpuUsage()
		}
	}()

	if len(config.Autostart) > 0 {
		startApplication(config.Autostart)
	}

	// Start server and listen for incoming requests
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.Port), nil)
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

func saveConfig(c *Config, filename string) error {
	if filename == "" {
		return errors.New("You should specify a filename to save the config")
	}

	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0755)
	if err != nil {
		return err
	}
	return nil
}
