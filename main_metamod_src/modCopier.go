package main

import (
	"os/exec"
	"path/filepath"
)

func copyMods(modList []string, modsDirectory string, SoW_directory string) {
	for _, mod := range modList {
		copyFolder(filepath.Join(modsDirectory, mod), SoW_directory)
	}
}

// completely arbitrary paths
func copyFolder(source string, destination string) {

	cmd := exec.Command("cp", "--recursive", source, destination)

	_ = cmd.Run()
}
