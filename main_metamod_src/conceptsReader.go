package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func readConceptFiles(conceptsDir string) map[string]*concept {
	concepts := make(map[string]*concept)
	fileInfo, err := ioutil.ReadDir(conceptsDir)
	if err != nil {
		fmt.Println("failed to read directory: " + conceptsDir)
		log.Fatal(err)
	}

	for i := 0; i < len(fileInfo); i++ {
		if strings.Contains(fileInfo[i].Name(), ".txt") {
			concept := readConceptFile(filepath.Join(conceptsDir, fileInfo[i].Name()))
			concepts[concept.name] = concept
		}
	}

	return concepts
}

type concept struct {
	name string
	instances map[string]*instance
}

// Return instance file contents
func readConceptFile(file string) *concept {
	thisFile, err := os.Open(file)
	if err != nil {
		fmt.Println("Failed to open file: " + file)
		log.Fatal(err)
	}
	// Initialize scanner
	scanner := bufio.NewScanner(thisFile)

	var concept concept

	concept.name = strings.Split(filepath.Base(file), ".")[0]
	concept.instances = make(map[string]*instance)

	for scanner.Scan() {
		// get next line
		line := cleanLine(scanner.Text())
		fields := strings.Fields(line)
		if len(fields) > 0 {
			var instance instance
			instance.name = fields[0]
			instance.params = make(map[string]string)
			for i := 1; i < len(fields); i++ {
				splitField := strings.Split(fields[i], ":")
				if len(splitField) > 1 {
					instance.params[splitField[0]] = splitField[1]
				}
			}
			concept.instances[instance.name] = &instance
		}
	}
	return &concept
}

type instance struct {
	name string
	params map[string]string
}




