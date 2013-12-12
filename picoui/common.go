package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

// Use a wrapper to differentiate logged panics from unexpected ones.
type LoggedError struct{ error }

type AppInfo struct {
	Name string
}

const (
	PICOUI_IMPORT_PATH = "github.com/hypebeast/picoui/picoui-lib"
)

func panicOnError(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Abort: %s: %s\n", msg, err)
		panic(LoggedError{err})
	}
}

func copyFile(destFilename string, srcFilename string) {
	destFile, err := os.Create(destFilename)
	panicOnError(err, "Failed to create file: "+destFilename)

	srcFile, err := os.Open(srcFilename)
	panicOnError(err, "Failed to open: "+srcFilename)

	_, err = io.Copy(destFile, srcFile)
	panicOnError(err, fmt.Sprintf("Failed to copy data from %s to %s\n", destFilename, srcFilename))

	err = destFile.Close()
	panicOnError(err, "Failed to close: "+destFilename)

	err = srcFile.Close()
	panicOnError(err, "Failed to close: "+srcFilename)
}

// copyDir copies a directory tree over to a new directory. Every file that ends
// with .template will be
func copyDir(destDir string, srcDir string, data interface{}) error {
	return filepath.Walk(srcDir, func(srcPath string, info os.FileInfo, err error) error {
		relSrcPath := strings.TrimLeft(srcPath[len(srcDir):], string(os.PathSeparator))
		destPath := filepath.Join(destDir, relSrcPath)

		// Skip dot files and dot directories
		if strings.HasPrefix(relSrcPath, ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Handle subdirectories
		if info.IsDir() {
			err := os.MkdirAll(destPath, 0755)
			if !os.IsExist(err) {
				panicOnError(err, "Failed to create directory")
			}
			return nil
		}

		// Handle templates: if a file ends in ".template", render it
		if strings.HasSuffix(relSrcPath, ".template") {
			renderTemplate(destPath[:len(destPath)-len(".template")], srcPath, data)
			return nil
		}

		// Just copy 'normal' files
		copyFile(destPath, srcPath)
		return nil
	})
}

func renderTemplate(destPath string, srcPath string, data interface{}) {
	tmpl, err := template.ParseFiles(srcPath)
	panicOnError(err, "Failed to parse template: "+srcPath)

	f, err := os.Create(destPath)
	panicOnError(err, "Failed to create file: "+destPath)

	err = tmpl.Execute(f, data)
	panicOnError(err, "Failed to render template: "+srcPath)

	err = f.Close()
	panicOnError(err, "Failed to close: "+f.Name())
}

// if os.env DEBUG set, debug is on
// Taken from: https://github.com/beego/bee/blob/master/util.go
func Debugf(format string, a ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "<unknown>"
			line = -1
		} else {
			file = filepath.Base(file)
		}
		fmt.Fprintf(os.Stderr, fmt.Sprintf("[debug] %s:%d %s\n", file, line, format), a...)
	}
}

func isPathInGopath(path string) bool {
	gopath := os.Getenv("GOPATH")
	isingopath := false
	splitedgopath := filepath.SplitList(gopath)
	for _, sp := range splitedgopath {
		sp, _ = filepath.EvalSymlinks(filepath.Join(sp, "src"))

		if strings.HasPrefix(strings.ToLower(path), strings.ToLower(sp)) {
			isingopath = true
			break
		}
	}
	return isingopath
}
