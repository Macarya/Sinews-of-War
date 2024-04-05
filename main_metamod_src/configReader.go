package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type configData struct {
	modOutFolder string

	mapDataFolder string
	mapDataBackupFolder string
	landedTitleFolder string

	innovationsFolder string
	tenetsFolder string
	traditionsFolder string
	pillarsFolder string

	terrainFolder string
	metamodFolder string

	onActionsFolder string

	processMap bool
	printEvents bool
	patchMode bool

	pixel_exponent float64

	dependencies []string

	// derivative parameters
	scriptedEffectsOutFolder string
	localizationOutFolder    string
	scriptedGUIOutFolder     string

	tradeNodePeriod int
	tradeNodeOffset int
}

func readConfigFile(config string) configData {

	thisFile, err := os.Open(config)
	if err != nil {
		log.Fatal(err)
	}
	dependentMods := make([]string,0)
	var configInfo configData
	// Initialize scanner
	scanner := bufio.NewScanner(thisFile)

	for scanner.Scan() {
		// get next line
		line := scanner.Text()
		fields := strings.Split(line, ";")
		if len(fields) > 1 {
			if fields[0] == "map_data_folder" {
				configInfo.mapDataFolder = fields[1]
			} else if fields[0] == "map_data_backup_folder" {
				configInfo.mapDataBackupFolder = fields[1]
			} else if fields[0] == "landed_titles_folder" {
				configInfo.landedTitleFolder = fields[1]
			} else if fields[0] == "province_terrain_folder" {
				configInfo.terrainFolder = fields[1]
			} else if fields[0] == "innovations_folder" {
				configInfo.innovationsFolder = fields[1]
			} else if fields[0] == "doctrines_folder" {
				configInfo.tenetsFolder = fields[1]
			} else if fields[0] == "pillars_folder" {
				configInfo.pillarsFolder = fields[1]
			} else if fields[0] == "traditions_folder" {
				configInfo.traditionsFolder = fields[1]
			} else if fields[0] == "output_mod_folder" {
				configInfo.modOutFolder = fields[1]
			} else if fields[0] == "metamod_folder" {
				configInfo.metamodFolder = fields[1]
			} else if fields[0] == "on_actions_folder" {
				configInfo.onActionsFolder = fields[1]
			} else if fields[0] == "pixel_exponent" {
				configInfo.pixel_exponent, err = strconv.ParseFloat(fields[1], 64)
				if err != nil {
					fmt.Println("Error setting parameter \"pixel_exponent\". \"" +  fields[1] + "\" is not a valid float64.")
					log.Fatal(err)
				}
			} else if fields[0] == "process_map" {
				if fields[1] == "true" || fields[1] == "yes" || fields[1] == "True" || fields[1] == "Yes" {
					configInfo.processMap = true
				}
			} else if fields[0] == "print_events" {
				if fields[1] == "true" || fields[1] == "yes" || fields[1] == "True" || fields[1] == "Yes" {
					configInfo.printEvents = true
				}
			} else if fields[0] == "patch_mode" {
				if fields[1] == "true" || fields[1] == "yes" || fields[1] == "True" || fields[1] == "Yes" {
					configInfo.patchMode = true
				}
			} else if fields[0] == "copy_mod" {
				dependentMods = append(dependentMods, fields[0])
			} else if fields[0] == "trade_node_period" {
				configInfo.tradeNodePeriod, err = strconv.Atoi(fields[1])
				if err != nil {
					fmt.Println("Nota valid int for trade_node_period")
				}
			} else if fields[0] == "trade_node_offset" {
				configInfo.tradeNodeOffset, err = strconv.Atoi(fields[1])
				if err != nil {
					fmt.Println("Nota valid int for trade_node_offset")
				}
			}
		}
	}
	if configInfo.tradeNodePeriod == 0 {
		configInfo.tradeNodePeriod = tradeNodePeriod
	}
	if configInfo.tradeNodeOffset == 0 {
		configInfo.tradeNodeOffset = tradeNodeOffset
	}
	// check for errors
	paramError := false
	if len(configInfo.mapDataFolder) < 1 {
		fmt.Println("Config File is missing parameter \"map_data_folder\"")
		paramError = true
	}
	if len(configInfo.landedTitleFolder) < 1 {
		fmt.Println("Config File is missing parameter \"landed_titles_folder\"")
		paramError = true
	}
	if len(configInfo.terrainFolder) < 1 {
		fmt.Println("Config File is missing parameter \"province_terrain_folder\"")
		paramError = true
	}
	if len(configInfo.innovationsFolder) < 1 {
		fmt.Println("Config File is missing parameter \"innovations_folder\"")
		paramError = true
	}
	if len(configInfo.metamodFolder) < 1 {
		fmt.Println("Config File is missing parameter \"metamod_folder\"")
		paramError = true
	}
	if len(configInfo.tenetsFolder) < 1 {
		fmt.Println("Config File is missing parameter \"doctrines_folder\"")
		paramError = true
	}
	if len(configInfo.traditionsFolder) < 1 {
		fmt.Println("Config File is missing parameter \"traditions_folder\"")
		paramError = true
	}
	if len(configInfo.pillarsFolder) < 1 {
		fmt.Println("Config File is missing parameter \"pillars_folder\"")
		paramError = true
	}
	if len(configInfo.onActionsFolder) < 1 {
		fmt.Println("Config File is missing parameter \"on_actions_folder\"")
		paramError = true
	}
	if configInfo.pixel_exponent == 0.0 {
		fmt.Println("Config File is missing parameter \"pixel_exponent\". Setting to default value of 3.")
	}
	if paramError {
		err := "Config file is missing parameters"
		log.Fatal(err)
	}
	if len(configInfo.modOutFolder) < 1 {
		fmt.Println("Config File is missing parameter \"output_mod_folder\"")
		paramError = true
	} else {
		configInfo.scriptedEffectsOutFolder = filepath.Join(configInfo.modOutFolder, "common", "scripted_effects")
		configInfo.localizationOutFolder = filepath.Join(configInfo.modOutFolder, "localization", "english")
		configInfo.scriptedGUIOutFolder = filepath.Join(configInfo.modOutFolder, "common", "scripted_guis")
	}
	return configInfo
}
