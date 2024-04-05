package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func readTemplateFiles(inDir string, configInfo configData) []string {
	fileInfo, err := ioutil.ReadDir(inDir)
	if err != nil {
		fmt.Println("failed to read directory: " + inDir)
		log.Fatal(err)
	}

	var tradeNodeStruggles []string
	for i := 0; i < len(fileInfo); i++ {
		newTemplate := readTemplateFile(inDir, fileInfo[i].Name(), configInfo)
		writeTemplateFile(newTemplate)
		if newTemplate.key == "DEMD_META_REGION" {
			tradeNodeStruggles = newTemplate.instances
			outFile, _ := os.Create(filepath.Join(configInfo.modOutFolder, "localization", "english", "demd_population_struggle_l_english.yml"))
			writeLocHeader(outFile)
			for _, instance := range newTemplate.instances {
				outFile.WriteString("demd_trade_node_struggle_" + instance + ":0 \"Trade Node: $" + instance + "$\"\n")
			}

		}
	}
	return tradeNodeStruggles
}

type template struct {
	key string
	instances []string
	outFolder string
	outName string
	lines []string
}

func writeTemplateFile(newTemplate template) {
	err := os.MkdirAll(newTemplate.outFolder, 0755)
	if err != nil {
		log.Fatal("failed to create folder: " + newTemplate.outFolder)
	}
	outFile, err := os.Create(filepath.Join(newTemplate.outFolder, newTemplate.outName))
	if err != nil {
		fmt.Println("Failed to create new file: " + filepath.Join(newTemplate.outFolder, newTemplate.outName))
		time.Sleep(100)
		log.Fatal(err)
	}
	writeHeader(outFile)
	for _, instance := range newTemplate.instances {
		for _, line := range newTemplate.lines {
			line = strings.ReplaceAll(line, newTemplate.key, instance)
			_, err = outFile.WriteString(line + "\n")
		}
		_, err = outFile.WriteString("\n\n")
	}

}

func readTemplateFile(inDir string, fileName string, configInfo configData) template {
	inFile := filepath.Join(inDir, fileName)

	thisFile, err := os.Open(inFile)
	if err != nil {
		fmt.Println("Failed to open file: " + inFile)
		log.Fatal(err)
	}
	// Initialize scanner
	scanner := bufio.NewScanner(thisFile)

	var newTemplate template
	newTemplate.outName = fileName

	for scanner.Scan() {
		line := cleanLine(scanner.Text())
		fields := strings.Fields(line)
		if len(fields) > 1 && fields[0] == ">>>" {
			newTemplate.key = fields[1]
		} else if len(fields) > 1 && fields[0] == ">>" {
			newTemplate.outFolder = filepath.Join(configInfo.modOutFolder, fields[1])
		} else if len(fields) > 1 && fields[0] == ">" {
			newTemplate.instances = append(newTemplate.instances, fields[1])
		} else {
			newTemplate.lines = append(newTemplate.lines, line)
		}
	}
	return newTemplate
}

