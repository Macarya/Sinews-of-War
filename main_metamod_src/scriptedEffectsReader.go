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

func readScriptedEffectsFiles(inDir string) map[string]*scriptedEffect {
	fileInfo, err := ioutil.ReadDir(inDir)
	if err != nil {
		fmt.Println("failed to read directory: " + inDir)
		log.Fatal(err)
	}

	scriptedEffectsLib := make(map[string]*scriptedEffect)

	for i := 0; i < len(fileInfo); i++ {
		effects := readScriptedEffectFile(filepath.Join(inDir,fileInfo[i].Name()))
		for _, effect := range effects {
			scriptedEffectsLib[effect.name] = effect
		}

	}
	return scriptedEffectsLib
}

// Return instance file contents
func readScriptedEffectFile(inFile string) map[string]*scriptedEffect {
	thisFile, err := os.Open(inFile)
	if err != nil {
		fmt.Println("Failed to open file: " + inFile)
		log.Fatal(err)
	}
	// Initialize scanner
	scanner := bufio.NewScanner(thisFile)

	effects:= make(map[string]*scriptedEffect)

	var effect scriptedEffect
	effect.lines = []string{}
	var currentEffect *scriptedEffect
	currentEffect = &effect
	bracketCounter := 0
	for scanner.Scan() {
		line := cleanLine(scanner.Text())
		if strings.Contains(line, "}") {
			bracketCounter--
		}
		if strings.Contains(line, "{") {
			bracketCounter++
		}
		fields := strings.Fields(line)
		if len(fields) > 0 {
			if getIndentLevel(line) == 0 {
				if fields[0] != "}" {
					var effect scriptedEffect
					effect.lines = []string{}
					effect.name = fields[0]
					currentEffect = &effect


				} else {
					effects[currentEffect.name] = currentEffect
					if bracketCounter != 0 {
						fmt.Println("error: mismatched brackets in effect " + effect.name + " in file " + inFile)
					}
				}
			} else {
				currentEffect.lines = append(currentEffect.lines, line)
			}
		}
	}
	//fmt.Println(effect.name)
	return effects
}

// removes comment blocks from a line (string)
func cleanLine(line string) string {
	// Look out for following byte signifying a comment
	commentByte := []byte("#")[0]
	// iterate through all bytes in line
	for i := 0; i < len(line); i++ {
		// if byte = comment byte, return line up until that byte
		if line[i] == commentByte {
			if i > 1 { return line[0:i] } else { return "" }
		}
	}
	// if no comments found, return whole line
	return line
}

type scriptedEffect struct {
	name string
	lines []string
}