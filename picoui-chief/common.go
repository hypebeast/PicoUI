package main

type AppInfo struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Config struct {
	AppsFolder string `json:"appFolder"`
	Port       int    `json:"port"`
}
