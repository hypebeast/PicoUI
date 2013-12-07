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
	"bufio"
	"bytes"
	"encoding/json"
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
	config      *Config
	app_running bool = false
	running_app *exec.Cmd
)

func readProc(file string) string {
	contents, err := ioutil.ReadFile(fmt.Sprintf("/proc/%s", file))
	if err != nil {
		return "error"
	}
	return string(contents[:])
}

func executeCmd(command string, args string) string {
	var out []byte
	var err error

	if args != "" {
		out, err = exec.Command(command, args).Output()
	} else {
		out, err = exec.Command(command).Output()
	}
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

func cpu() []uint64 {
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

func startAppHandler(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	appName := vals["name"][0]

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
	running_app.Dir = foundApp.Path
	err := running_app.Start()
	if err != nil {
		fmt.Fprintf(w, "%s", "error: can't start app")
		return
	}

	app_running = true

	fmt.Fprintf(w, "%s:%s", "ok starting app: ", foundApp.Name)
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
	stats := cpu()
	fmt.Fprintf(w, "%d %d %d %d %d %d %d %d", stats[0], stats[1], stats[2], stats[3], stats[4], stats[5], stats[6], stats[7])
}

func cpuUtilHandler(w http.ResponseWriter, r *http.Request) {
	s1 := cpu()
	time.Sleep(100 * time.Millisecond)
	s2 := cpu()

	var total1 uint64 = 0
	var total2 uint64 = 0
	for i, _ := range s1 {
		total1 += s1[i]
		total2 += s2[i]
	}

	idle1 := s1[4]
	idle2 := s2[4]
	diff_idle := idle2 - idle1
	diff_total := total2 - total1
	diff_usage := (1000*(diff_total-diff_idle)/diff_total + 5) / 10

	//loadavg := math.Float64frombits(((s2[0] + s2[1] + s2[2]) - (s1[0] + s1[1] + s1[2])) / ((s2[0] + s2[1] + s2[2] + s2[3]) - (s1[0] + s1[1] + s1[2] + s1[3])))
	// fmt.Println(loadavg)
	fmt.Fprintf(w, "%d", diff_usage)
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
	http.HandleFunc("/apps", listAppsHandler)
	http.HandleFunc("/apps/start", startAppHandler)
	http.HandleFunc("/apps/kill", killAppHandler)
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
	http.HandleFunc("/system/cpu_util", cpuUtilHandler)
	http.HandleFunc("/ping", pingHandler)

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
