package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func scriptedEffectsPrinter(scriptedEffects map[string]*scriptedEffect, topLevelEffects []string, outDir string, outName string) {
	err := os.MkdirAll(filepath.Dir(filepath.Join(outDir,outName)), 0755)
	if err != nil {
		fmt.Println("Failed to create new directory: " + outDir)
		time.Sleep(100)
		log.Fatal(err)
	}
	outFile, err := os.Create(filepath.Join(outDir, outName))
	if err != nil {
		fmt.Println("Failed to create new file: " + filepath.Join(outDir, outName))
		time.Sleep(100)
		log.Fatal(err)
	}
	writeHeader(outFile)
	for _, effect := range scriptedEffects {
		if stringInList(effect.name, topLevelEffects) {
			_, err = outFile.WriteString(effect.name + " = {\n")
			for _, line := range effect.lines {
				_, err = outFile.WriteString(line + "\n")
			}
			_, err = outFile.WriteString("}\n\n")
		}
	}
}

func writeHeader(outFile *os.File) {
	_, _ = outFile.WriteString("\ufeff")
	_, _ = outFile.WriteString("#############################################\n")
	_, _ = outFile.WriteString("# DEMD Population System\n")
	_, _ = outFile.WriteString("# by Vertimnus\n")
	_, _ = outFile.WriteString("# This file was compiled by a machine from jomini metascript source code.\n")
	_, _ = outFile.WriteString("# It should never be manually edited.\n")
	_, _ = outFile.WriteString("#############################################\n\n")
}

func stringInList(string string, list []string) bool {
	for _, a := range list {
		if a == string {
			return true
		}
	}
	return false
}

func stringInListCount(string string, list []string) int {
	val := 0
	for _, a := range list {
		if a == string {
			val++
		}
	}
	return val
}