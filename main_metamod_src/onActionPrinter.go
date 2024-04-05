package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func onActionPrinter(configInfo configData, outName string, tradeNodeKeys []string) {
	err := os.MkdirAll(filepath.Dir(filepath.Join(configInfo.onActionsFolder,outName)), 0755)
	if err != nil {
		fmt.Println("Failed to create new directory: " + configInfo.onActionsFolder)
		time.Sleep(100)
		log.Fatal(err)
	}
	outFile, err := os.Create(filepath.Join(configInfo.onActionsFolder, outName))
	if err != nil {
		fmt.Println("Failed to create new file: " + filepath.Join(configInfo.onActionsFolder, outName))
		time.Sleep(100)
		log.Fatal(err)
	}
	writeHeader(outFile)

	for i := 1; i <= 365; i++ {
		_, err = outFile.WriteString("yearly_global_pulse = {\n")
		_, err = outFile.WriteString("\ton_actions = {\n")
		if i > 1 {
			_, err = outFile.WriteString("\t\tdelay = { days = " + strconv.Itoa(i-1) + " }\n")
		}
		_, err = outFile.WriteString("\t\tdemd_meta_population_" + strconv.Itoa(i) + "\n")
		_, err = outFile.WriteString("\t}\n")
		_, err = outFile.WriteString("}\n")
	}

	j := 1
	for i := 1; i <= 365; i++ {
		_, err = outFile.WriteString("demd_meta_population_" + strconv.Itoa(i) + " = {\n")
		_, err = outFile.WriteString("\teffect = {\n")

		_, err = outFile.WriteString("\t\tannualPulse = { list = eligible_counties_group_y_" + strconv.Itoa(i) + " }\n")
		if i < 361 {
			//_, err = outFile.WriteString("\t\t# monthlyPulse = { list = eligible_counties_group_m_" + strconv.Itoa(j) + " }\n")
			if j < len(tradeNodeKeys) + 1 { // +1 since trade keys go 1-30 to match days of month/year
				_, err = outFile.WriteString("\t\tmonthlyPulseTradeNode = { list = trade_nodes_group_m_" + strconv.Itoa(j) + " }\n")
			}

			j++
			if j > 30 {
				j = 1
			}
		}
		_, err = outFile.WriteString("\t}\n")
		_, err = outFile.WriteString("}\n")
	}
}
