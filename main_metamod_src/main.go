package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var tradeNodePeriod = 1
var tradeNodeOffset = 0

func main() {
	fmt.Println()
	fmt.Println("#################################################")
	fmt.Println("Sinews of War Auto-Patcher")
	fmt.Println("By Vertimnus")
	currentTime := time.Now()
	fmt.Println("Compiled " + currentTime.Month().String() + " " + strconv.Itoa(currentTime.Day()) + ", " + strconv.Itoa(currentTime.Year()))
	fmt.Println("#################################################")
	fmt.Println()
	time.Sleep(100)

	if len(os.Args) < 2 {
		fmt.Println("Please provide the path to a configuration file as the first argument.")
	} else {

		// read config file and retrieve params
		configInfo := readConfigFile(os.Args[1])

		bakEliminator(configInfo.modOutFolder)
		conceptsFolder := filepath.Join(configInfo.metamodFolder, "concepts")
		scriptedEffectsFolder := filepath.Join(configInfo.metamodFolder, "scriptedEffects")
		templatesFolder := filepath.Join(configInfo.metamodFolder, "templates")



		// parse concepts and scripted effects
		fmt.Println("reading concepts...")
		concepts := readConceptFiles(conceptsFolder)
		fmt.Println("reading scripted effects...")
		scriptedEffectsLib := readScriptedEffectsFiles(scriptedEffectsFolder)
		fmt.Println("expanding metacode...")
		scriptedEffectsLib = metaCodeExpander(scriptedEffectsLib, concepts)
		fmt.Println("expanding scripted effects...")
		scriptedEffectsLib = scriptedEffectsExpander(scriptedEffectsLib, concepts)

		// delete temporary meta files
		/*err := os.Remove(filepath.Join(scriptedEffectsFolder, "temp_innovations.txt"))
		if err != nil {
			log.Fatal(err)
		}
		err = os.Remove(filepath.Join(scriptedEffectsFolder, "temp_tenets.txt"))
		if err != nil {
			log.Fatal(err)
		}
		err = os.Remove(filepath.Join(scriptedEffectsFolder, "temp_pillars.txt"))
		if err != nil {
			log.Fatal(err)
		}
		err = os.Remove(filepath.Join(scriptedEffectsFolder, "temp_traditions.txt"))
		if err != nil {
			log.Fatal(err)
		}*/

		// output things
		if !configInfo.patchMode {
			// output stuff read from innos/tenets in the metamod scripted effect folder
			fmt.Println("writing inno meta-meta...")
			writeInnoScriptedEffect(configInfo.innovationsFolder, scriptedEffectsFolder, "temp_innovations.txt")
			tenetLib := writeTenetScriptedEffect(configInfo.tenetsFolder, scriptedEffectsFolder, "temp_tenets.txt")
			tradLib := writeTraditionScriptedEffect(configInfo.traditionsFolder, scriptedEffectsFolder, "temp_traditions.txt")
			writePillarScriptedEffect(configInfo.pillarsFolder, scriptedEffectsFolder, "temp_pillars.txt")

			fmt.Println("Writing Template Files...")
			tradeNodeKeys := readTemplateFiles(templatesFolder, configInfo)

			writeAuxiliaryScriptedEffects(configInfo, tradLib, tenetLib, tradeNodeKeys)

			fmt.Println("Writing CCU Files...")
			writeCCUFiles(configInfo)



			fmt.Println("writing annual amortized pulse...")
			scriptedEffectsPrinter(scriptedEffectsLib, []string{"annualPulse", "monthlyPulseTradeNode"}, configInfo.scriptedEffectsOutFolder, "demd_population_annual_pulse.txt")
			fmt.Println("writing monthly amortized pulse...")
			scriptedEffectsPrinter(scriptedEffectsLib, []string{"monthlyPulse"}, configInfo.scriptedEffectsOutFolder, "demd_population_monthly_pulse.txt")
			fmt.Println("writing annual un-amortized pulse...")
			scriptedEffectsPrinter(scriptedEffectsLib, []string{"trade_node_pulse", "fertilitySubRoutine", "ruler_pulse", "world_migration_end",
				"world_migration_start"}, configInfo.scriptedEffectsOutFolder, "demd_population_year_end_pulse.txt")
			fmt.Println("writing startup pulse...")
			scriptedEffectsPrinter(scriptedEffectsLib, []string{"startupPulse"}, configInfo.scriptedEffectsOutFolder, "demd_population_startup_pulse.txt")
			fmt.Println("writing misc functions...")
			scriptedEffectsPrinter(scriptedEffectsLib, []string{"setEdicts", "empire_title_pulse", "faithEconomy", "cultureEconomy", "apply_province_garrison", "apply_province_levy"}, configInfo.scriptedEffectsOutFolder, "demd_population_single_province_economy.txt")

			fmt.Println("writing on actions and events...")
			onActionPrinter(configInfo, "demd_meta_on_actions.txt", tradeNodeKeys)
		}


		if configInfo.processMap {
			fmt.Println("initializing provs...")
			landedTitlePath := filepath.Join(configInfo.landedTitleFolder, "00_landed_titles.txt")
			provinceWinterFile := filepath.Join(configInfo.terrainFolder, "01_province_properties.txt")
			provinceFunction(configInfo.mapDataFolder, configInfo.mapDataBackupFolder, provinceWinterFile, landedTitlePath, configInfo.scriptedEffectsOutFolder, "demd_initializer_effects.txt", configInfo.pixel_exponent)
		}

		fmt.Println("Copying dependencies...")
		path, _ := os.Getwd()
		SoWdirectory := filepath.Dir(path)
		modsDirectory := filepath.Dir(SoWdirectory)
		copyMods(configInfo.dependencies, modsDirectory, SoWdirectory)


	}
}

func bakEliminator(path string) {
	_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf(err.Error())
		}
		if !info.IsDir() && len(strings.Split(info.Name(), ".")) > 1 && strings.Split(info.Name(), ".")[len(strings.Split(info.Name(), "."))-1] == "bak" {
			fmt.Printf("Deleting: %s\n", info.Name())
			err = os.Remove(path)
			if err != nil {
				log.Fatal(err)
			}
		}

		return nil
	})
}