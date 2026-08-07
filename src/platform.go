package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func rootPath() (string, string) {
	var root string
	if runtime.GOOS == "windows" {
		root = filepath.Join(os.Getenv("APPDATA"), "Adobe", "CoreSync", "plugins", "livetype")
	} else if runtime.GOOS == "darwin" {
		home, _ := os.UserHomeDir()
		root = filepath.Join(home, "Library", "Application Support", "Adobe", "CoreSync", "plugins", "livetype")
	}
	if root != "" {
		if _, _, err := findDirectories(root); err == nil {
			return root, ""
		}
	}
	if _, _, err := findDirectories("."); err == nil {
		return ".", ""
	}
	if _, _, err := findDirectories("livetype"); err == nil {
		return "livetype", ""
	}
	return "", "Adobe Fonts were not found. Adobe may not be installed, or you may be running Linux. Choose the source folder manually."
}

func uriPath(path string) string {
	path = filepath.FromSlash(path)
	if runtime.GOOS == "windows" && strings.HasPrefix(path, `\`) && len(path) > 2 && path[2] == ':' {
		return path[1:]
	}
	return path
}

func openDirectory(path string) error {
	command := "xdg-open"
	if runtime.GOOS == "windows" {
		command = "explorer.exe"
	} else if runtime.GOOS == "darwin" {
		command = "open"
	}
	return exec.Command(command, path).Start()
}
