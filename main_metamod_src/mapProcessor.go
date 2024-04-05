package main

import (
	"bufio"
	"fmt"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func provinceFunction(mapDataFolder string, mapDataBackupFolder string, provinceWinterFile string, landedTitleFile string, outDir string, outName string, pixelExp float64) {

	fmt.Println("processing definition.csv")
	provinces, color2provID := getProvinces(mapDataFolder, mapDataBackupFolder)
	fmt.Println("processing provinces.png")
	getProvinceInfo(mapDataFolder, mapDataBackupFolder, provinces, color2provID)
	fmt.Println("processing 00_landed_titles.txt")
	_ = getTitleInfo(provinces, landedTitleFile)
	fmt.Println("computing province positions and elevation")
	getAvgPos(provinces)
	fmt.Println("finding major river and strait crossings")
	getAdjacencies(provinces, mapDataFolder, mapDataBackupFolder)
	fmt.Println("finding and classifying water provinces")
	getWater(provinces)
	fmt.Println("finding neighboring and nearby provinces")
	getNeighbors(provinces, mapDataFolder, mapDataBackupFolder, color2provID)
	getNearby(provinces)
	fmt.Println("finding provinces on minor rivers")
	getMinorRivers(provinces, mapDataFolder, mapDataBackupFolder)
	fmt.Println("finding provinces on major bodies of water")
	getMajorWaterAdjacency(provinces)
	fmt.Println("retrieving winter data")
	getWinterData(provinceWinterFile,provinces)
	fmt.Println("converting province data to county data")
	counties := province2county(provinces)
	fmt.Println("writing to disk")
	writeProvinces(counties, provinces, outDir, outName, pixelExp)
}

func writeProvinces(counties map[string]*county, provinces map[string]*prov, outDir string, outName string, pixelExp float64) {
	err := os.MkdirAll(filepath.Dir(filepath.Join(outDir,outName)), 0755)
	outFile, err := os.Create(filepath.Join(outDir, outName))
	if err != nil {
		fmt.Println("Failed to create new file: " + filepath.Join(outDir, outName))
		time.Sleep(100)
		log.Fatal(err)
	}
	fmt.Println("created file at " + filepath.Join(outDir, outName))

	writeHeader(outFile)
	_, _ = outFile.WriteString("initialize_province_variables = {\n")
	_, _ = outFile.WriteString("\n")

	counter := 1
	counterM := 1
	for _, province := range counties {
		if !(province.x == 0 && province.y == 0) && !math.IsNaN(province.x) && !math.IsNaN(province.y) { // try to weed out bad provinces from TFE
			// monthly split
			_, _ = outFile.WriteString("\t# add_to_global_variable_list = { name = eligible_counties_group_m_" + strconv.Itoa(counterM) + " target = title:" + province.name + " }\n")
			// annual split
			_, _ = outFile.WriteString("\tadd_to_global_variable_list = { name = eligible_counties_group_y_" + strconv.Itoa(counter) + " target = title:" + province.name + " }\n")
			// annual split
			_, _ = outFile.WriteString("\tadd_to_global_variable_list = { name = eligible_counties_all target = title:" + province.name + " }\n")
			counter++
			counterM++
			if counter > 365 {
				counter = 1
			}
			if counterM > 30 {
				counterM = 1
			}
		}
	}

	_, _ = outFile.WriteString("\n")
	for _, province := range counties {
		_, _ = outFile.WriteString("\ttitle:" + province.name + " = {\n")
		_, _ = outFile.WriteString("\t\tset_variable = { name = id value = " + province.capitalID + " }\n")
		/*if !math.IsNaN(province.x) && !math.IsNaN(province.y) { // print placeholder vals if x&y are NaN. Shouldn't be necessary but when you work with Wrench...
			_, _ = outFile.WriteString("\t\tset_variable = { name = x value = " + fmt.Sprintf("%.3f", province.x) + " }\n")
			_, _ = outFile.WriteString("\t\tset_variable = { name = y value = " + fmt.Sprintf("%.3f", province.y) + " }\n")
		} else {
			_, _ = outFile.WriteString("\t\tset_variable = { name = x value = 0 }\n")
			_, _ = outFile.WriteString("\t\tset_variable = { name = y value = 0 }\n")
		}*/
		_, _ = outFile.WriteString("\t}\n\n")
	}

	// write standard info
	for id, province := range provinces {
		if province.isValidForHolding {
			_, _ = outFile.WriteString("\tprovince:" + id + " = { # " + province.name + "\n")
			if province.isOnLake {
				_, _ = outFile.WriteString("\t\tset_variable = isOnLake\n")
				_, _ = outFile.WriteString("\t\tcounty = { set_variable = isOnLake }\n")
			}
			/*if province.isOnMajorRiver {
				_, _ = outFile.WriteString("\t\tset_variable = { name = isOnMajorRiver value = 1 }\n")
				_, _ = outFile.WriteString("\t\tcounty = { set_variable = { name = isOnMajorRiver value = 1 } }\n")
			}*/
			if province.isOnMinorRiver {
				_, _ = outFile.WriteString("\t\tset_variable = isOnMinorRiver\n")
				_, _ = outFile.WriteString("\t\tcounty = { set_variable = isOnMinorRiver }\n")
			}
			if province.isMajorRiverCrossing {
				_, _ = outFile.WriteString("\t\tset_variable = isMajorRiverCrossing\n")
				_, _ = outFile.WriteString("\t\tcounty = { set_variable = isMajorRiverCrossing }\n")
			}
			if province.isStraitCrossing {
				_, _ = outFile.WriteString("\t\tset_variable = isStraitCrossing\n")
				_, _ = outFile.WriteString("\t\tcounty = { set_variable = isStraitCrossing }\n")
			}
			_, _ = outFile.WriteString("\t\tset_variable = { name = winter_severity_bias value = " + fmt.Sprintf("%.2f", province.winterSeverityBias) + " }\n")
			_, _ = outFile.WriteString("\t\tset_variable = { name = sqrt_pixels value = " + fmt.Sprintf("%.2f", math.Pow(float64(len(province.pixels)), 1.0/pixelExp)) + " } # sqrt pixels " +
			 fmt.Sprintf("%.2f", math.Pow(float64(len(province.pixels)), 0.5)) + "\n")

			_, _ = outFile.WriteString("\t}\n\n")
		}
	}
	/*_, _ = outFile.WriteString("\t# Water Adjacencies\n\n")
	for id, province := range provinces {
		if !strings.Contains(province.ID, "#") && (province.isSea || province.isMajorRiver || province.isLake) {
			_, _ = outFile.WriteString("\tprovince:" + id + " = { # " + province.name + "\n")
			for _, neighborID := range province.neighbors {

				neighbor := provinces[neighborID]
				if neighbor.isValidForHolding && neighbor.isCountyCapital {
					neighborCounty := counties[neighbor.county]
					_, _ = outFile.WriteString("\t\ttitle:" + neighborCounty.name + " = {\n")
					_, _ = outFile.WriteString("\t\t\tprev = { add_to_variable_list = { name = adjacent_counties target = prev } }")
					_, _ = outFile.WriteString("\t\t}\n")
				}
			}
			_, _ = outFile.WriteString("\t}\n")
		}
	}*/
	_, _ = outFile.WriteString("}")
}

type prov struct {
	name      string
	x         float64
	y         float64
	elevation float64
	terrain   string
	geoRegion string

	color     string
	ID        string
	neighbors []string
	nearby    []string
	port      string
	pixels    []*pixel

	baseWater     float64
	baseFertility float64

	isValidForHolding    bool
	isLand               bool
	isOnMinorRiver       bool
	isOnMajorRiver       bool
	isOnLake             bool
	isMajorRiverCrossing bool
	isStraitCrossing     bool
	isCoastal            bool

	isLake       bool
	isMajorRiver bool
	isSea        bool

	barony  string
	county  string
	duchy   string
	kingdom string
	empire  string

	isDuchyCapital  bool
	isCountyCapital bool

	duchyCapital  string
	countyCapital string

	duchyProvinces  []string
	countyProvinces []string

	pixelArea          float64
	farmDistrictBase   float64
	transportCostModifier   float64
	winterSeverityBias float64

	tradePower float64
}

type county struct {
	name string
	x    float64
	y    float64

	elevation float64
	latitude  float64

	geoRegion string

	capitalID string

	provIDs   []string // does not include capital
	neighbors []string

	terrains     []string
	terrainFracs map[string]float64

	coastalFrac    float64
	lakeFrac       float64
	minorRiverFrac float64
	majorRiverFrac float64
	resource       string
	resourceBonus float64
	resourceBonusMax float64

	terrainSanitation float64

	isDuchyCapital bool

	farmDistrictsBase float64
	fertilityBase     float64

	defenseShareMajor float64
	defenseShareMinor float64
}

func province2county(provinces map[string]*prov) map[string]*county {
	counties := make(map[string]*county)
	for _, province := range provinces {
		if province.isCountyCapital {
			var newCounty county
			newCounty.name = province.county
			newCounty.x = province.x
			newCounty.y = province.y
			newCounty.elevation = province.elevation
			newCounty.geoRegion = province.geoRegion
			newCounty.capitalID = province.ID
			newCounty.provIDs = province.countyProvinces
			newCounty.farmDistrictsBase = province.farmDistrictBase
			newCounty.fertilityBase = province.baseFertility
			for _, prv := range province.countyProvinces {
				newCounty.fertilityBase += provinces[prv].baseFertility
			}
			newCounty.fertilityBase = newCounty.fertilityBase / float64(len(province.countyProvinces)+1)


			newCounty.terrains = []string{}
			newCounty.isDuchyCapital = province.isDuchyCapital

			if province.isCoastal {
				newCounty.coastalFrac = 1
			} else {
				newCounty.coastalFrac = 0
			}
			for _, countyProvID := range newCounty.provIDs {
				if provinces[countyProvID].isCoastal {
					newCounty.coastalFrac++
				}
			}
			newCounty.coastalFrac = newCounty.coastalFrac / float64(len(newCounty.provIDs) + 1)

			if province.isOnLake {
				newCounty.lakeFrac = 1
			} else {
				newCounty.lakeFrac = 0
			}
			for _, countyProvID := range newCounty.provIDs {
				if provinces[countyProvID].isOnLake {
					newCounty.lakeFrac++
				}
			}
			newCounty.lakeFrac = newCounty.lakeFrac / float64(len(newCounty.provIDs) + 1)

			if province.isOnMinorRiver {
				newCounty.minorRiverFrac = 1
			} else {
				newCounty.minorRiverFrac = 0
			}
			for _, countyProvID := range newCounty.provIDs {
				if provinces[countyProvID].isOnMinorRiver {
					newCounty.minorRiverFrac++
				}
			}
			newCounty.minorRiverFrac = newCounty.minorRiverFrac / float64(len(newCounty.provIDs) + 1)

			if province.isOnMajorRiver {
				newCounty.majorRiverFrac = 1
			} else {
				newCounty.majorRiverFrac = 0
			}
			for _, countyProvID := range newCounty.provIDs {
				if provinces[countyProvID].isOnMajorRiver {
					newCounty.majorRiverFrac++
				}
			}
			newCounty.majorRiverFrac = newCounty.majorRiverFrac / float64(len(newCounty.provIDs) + 1)

			// sum up farming cap
			for _, countyProvID := range newCounty.provIDs {
				newCounty.farmDistrictsBase += provinces[countyProvID].farmDistrictBase
				newCounty.terrains = append(newCounty.terrains, provinces[countyProvID].terrain)
				//newCounty.fertilityBase += provinces[countyProvID].baseFertility
			}

			counties[newCounty.name] = &newCounty
		}

	}

	return counties

}

type title struct {
	name string
	tier int
	geoRegion string
	capitalProvince string
	capitalBarony string
	capitalCounty string
	titleAbove *title
	titleBelow []*title
	provinces []string
}

func getCleanedString(string string) string {
	cleanString := strings.ReplaceAll(string, "=", "")
	return cleanString
}

func getTitleInfo(provinces map[string]*prov, landedTitlesPath string) map[string]*title {
	// open file
	thisFile, err := os.Open(landedTitlesPath)
	if err != nil {
		fmt.Println("Failed to open landed titles file: " + landedTitlesPath)
		log.Fatal(err)
	}
	// Initialize scanner
	scanner := bufio.NewScanner(thisFile)

	titles := make(map[string]*title)

	var currentTitleReal title
	currentTitleReal.tier = -1
	currentTitleReal.provinces = []string{}
	currentTitleReal.titleBelow = []*title{}
	currentEmpire := &currentTitleReal
	currentKingdom := &currentTitleReal
	currentDuchy := &currentTitleReal
	currentCounty := &currentTitleReal
	currentBarony := &currentTitleReal
	currentTitle := &currentTitleReal
	baroniesRead := 0
	for scanner.Scan() {
		// get next line
		line := scanner.Text()
		line = cleanLine(line)
		fields := strings.Fields(line)
		if len(fields) > 0 {
			identifier := strings.Split(fields[0], "_")[0]
			if identifier == "e" {
				var newEmpire title
				newEmpire.name = fields[0]
				newEmpire.tier = 4
				currentEmpire = &newEmpire
				currentTitle = &newEmpire
				titles[newEmpire.name] = &newEmpire
			} else if identifier == "k" {
				var newKingdom title
				newKingdom.name = fields[0]
				newKingdom.tier = 3
				newKingdom.titleAbove = currentEmpire
				currentEmpire.titleBelow = append(currentEmpire.titleBelow, &newKingdom)
				currentKingdom = &newKingdom
				currentTitle = &newKingdom
				titles[newKingdom.name] = &newKingdom
			} else if identifier == "d" {
				var newDuchy title
				newDuchy.name = fields[0]
				newDuchy.tier = 2
				newDuchy.titleAbove = currentKingdom
				currentKingdom.titleBelow = append(currentKingdom.titleBelow, &newDuchy)
				currentDuchy = &newDuchy
				currentTitle = &newDuchy
				titles[newDuchy.name] = &newDuchy
			} else if identifier == "c" {
				var newCounty title
				newCounty.name = getCleanedString(fields[0])
				newCounty.tier = 1
				newCounty.titleAbove = currentDuchy
				newCounty.capitalCounty = newCounty.name
				currentDuchy.titleBelow = append(currentDuchy.titleBelow, &newCounty)
				currentCounty = &newCounty
				currentTitle = &newCounty
				titles[newCounty.name] = &newCounty
				baroniesRead = 0
			} else if identifier == "b" {
				var newBarony title
				newBarony.name = fields[0]
				newBarony.tier = 0
				newBarony.titleAbove = currentCounty
				newBarony.capitalBarony = newBarony.name
				currentCounty.titleBelow = append(currentCounty.titleBelow, &newBarony)
				currentBarony = &newBarony
				currentTitle = &newBarony
				titles[newBarony.name] = &newBarony
				if baroniesRead == 0 {
					currentCounty.capitalBarony = newBarony.name
				}
				baroniesRead++
			} else if fields[0] == "province" && len(fields) > 2 && currentTitle.tier == 0 {
				currentBarony.capitalProvince = fields[2]
				provinces[fields[2]].barony = currentBarony.name
				provinces[fields[2]].county = currentCounty.name
				provinces[fields[2]].duchy = currentDuchy.name
				provinces[fields[2]].kingdom = currentKingdom.name
				provinces[fields[2]].empire = currentEmpire.name
				provinces[fields[2]].isValidForHolding = true
				currentBarony.provinces = append(currentBarony.provinces, fields[2])
				currentCounty.provinces = append(currentCounty.provinces, fields[2])
				currentDuchy.provinces = append(currentDuchy.provinces, fields[2])
				currentKingdom.provinces = append(currentKingdom.provinces, fields[2])
				currentEmpire.provinces = append(currentEmpire.provinces, fields[2])
			} else if fields[0] == "capital" && len(fields) > 2 && currentTitle.tier > 1 {
				currentTitle.capitalCounty = fields[2]
			}
		}
	}

	// get the capital of the title as a province id rather than a county name (e.g. 482 rather than c_rome)
	for _, title := range titles {
		if title.tier > 1 {
			// capital = c_whatever to capital = b_whatever
			if capitalCounty, ok := titles[title.capitalCounty]; ok {
				title.capitalBarony = capitalCounty.capitalBarony
			} else {
				// fmt.Println("Failed to find capital county " + title.capitalCounty + " for title " + title.name)
			}
		}
		if title.tier > 0 {
			// capital = b_whatever to capital = prov id
			if capitalBarony, ok := titles[title.capitalBarony]; ok {
				title.capitalProvince = capitalBarony.capitalProvince
			} else {
				// fmt.Println("Failed to find capital barony " + title.capitalBarony + " for county " + title.name)
			}
		}
	}

	// identify county and duchy capital provinces
	for _, title := range titles {
		if title.tier == 1 {
			if capProv, ok := provinces[title.capitalProvince]; ok {
				capProv.isCountyCapital = true
			}
		} else if title.tier == 2 {
			if capProv, ok := provinces[title.capitalProvince]; ok {
				capProv.isDuchyCapital = true
			}
		}
	}
	// attach lists of all other provinces in the duchy/county to every duchy/county capital
	for _, province := range provinces {
		if province.isValidForHolding {
			for _, provInDuchy := range titles[province.duchy].provinces {
				if provInDuchy != province.ID {
					province.duchyProvinces = append(province.duchyProvinces, provInDuchy)
				}
			}
			for _, provInCounty := range titles[province.county].provinces {
				if provInCounty != province.ID {
					province.countyProvinces = append(province.countyProvinces, provInCounty)
				}
			}
			province.duchyCapital = titles[province.duchy].capitalProvince
			province.countyCapital = titles[province.county].capitalProvince
		}
	}
	return titles
}

func getNearby(provinces map[string]*prov) {
	for _, provA := range provinces {
		for _, provB := range provinces {
			if provA.ID < provB.ID && provA.isValidForHolding && provB.isValidForHolding {
				if getDistance(provA, provB) < 100.0 {
					provA.nearby = append(provA.nearby, provB.ID)
					provB.nearby = append(provB.nearby, provA.ID)
				}
			}
		}
	}
}

func getDistance(provA *prov, provB *prov) float64 {
	return math.Sqrt(math.Pow(provA.x - provB.x, 2) + math.Pow(provA.y - provB.y, 2))
}

func getWinterData(file string, provinces map[string]*prov) {
	thisFile, err := os.Open(file)
	if err != nil {
		fmt.Println("Failed to open file: " + file)
		log.Fatal(err)
	}
	// Initialize scanner
	scanner := bufio.NewScanner(thisFile)
	currentProv := ""
	isInsideProv := false
	for scanner.Scan() {
		// get next line

		line := cleanLine(scanner.Text())
		fields := strings.Fields(line)
		if len(fields) > 0 {
			if _, ok := provinces[fields[0]]; ok {
				currentProv = fields[0]
				isInsideProv = true
			} else if fields[0] == "winter_severity_bias" && isInsideProv && len(fields) > 2 {
				provinces[currentProv].winterSeverityBias, _ = strconv.ParseFloat(fields[2], 64)
			}
		}
	}
}

func getMajorWaterAdjacency(provinces map[string]*prov) {
	for _, province := range provinces {
		if province.isLand {
			for _, neighborID := range province.neighbors {
				if provinces[neighborID].isMajorRiver {
					province.isOnMajorRiver = true
				}
				if provinces[neighborID].isLake {
					province.isOnLake = true
				}
				if provinces[neighborID].isSea {
					province.isCoastal = true
				}
			}
		}
	}
}

func getWater(provinces map[string]*prov) {
	for _, province := range provinces {
		fractionLand := 0.0
		for _, pixel := range province.pixels {
			if pixel.class == "land" {
				fractionLand++
			}
		}
		fractionLand = fractionLand / float64(len(province.pixels))
		if fractionLand > 0.5 {
			province.isLand = true
		}
	}

	for _, province := range provinces {
		if !province.isLand {
			riverIdentifiers := []string{"river", "ilinalta", "honrich", "geir", "yorgrim"}
			lakeIdentifiers := []string{"lakes", "lake", "lagoon", "lagoons"}
			if nameIndicatesType(province.name, riverIdentifiers) {
				province.isMajorRiver = true
			} else if nameIndicatesType(province.name, lakeIdentifiers) {
				province.isLake = true
			} else {
				province.isSea = true
			}
		}
	}
}

func nameIndicatesType(name string, identifiers []string) bool {
	fields := strings.Fields(name)
	for _, field := range fields {
		for _, identifier := range identifiers {
			if strings.EqualFold(field, identifier) {
				return true
			}
		}
	}
	substrings := strings.Split(name, "_")
	for _, substring := range substrings {
		for _, identifier := range identifiers {
			if strings.EqualFold(substring, identifier) {
				return true
			}
		}
	}
	return false
}



func getMinorRivers(provinces map[string]*prov, mapDataFolder string, mapDataBackupFolder string) {
	// Open Province File
	riverFile, err := os.Open(filepath.Join(mapDataFolder, "rivers.png"))
	if err != nil {
		if len(mapDataBackupFolder) > 0 {
			riverFile, err = os.Open(filepath.Join(mapDataBackupFolder, "rivers.png"))
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println("Did not find rivers.png in \"" + mapDataFolder + "\"")
				fmt.Println("Using backup copy in \"" + mapDataBackupFolder + "\"")
			}
		} else {
			log.Fatal(err)
		}
	}

	riverMap, err := png.Decode(riverFile)
	if err != nil {
		log.Fatal(err)
	}

	for _, province := range provinces {
		for _, pix := range province.pixels {
			if pix.class == "river" {
				province.isOnMinorRiver = true
				break
			} else {
				neighborPixels := [][]int{{int(pix.x) + 1, int(pix.y) + 1}, {int(pix.x) - 1, int(pix.y) + 1}, {int(pix.x) + 1, int(pix.y) - 1}, {int(pix.x) - 1, int(pix.y) - 1}}
				for _, neighbor := range neighborPixels {
					r, g, b, _ := riverMap.At(neighbor[0], neighbor[1]).RGBA()
					unifiedColor := rgbaToUnified(int(r / 0x101), int(g / 0x101) , int(b / 0x101))
					if unifiedColor != "255.255.255" && unifiedColor != "255.0.128" {
						province.isOnMinorRiver = true
						break
					}
				}
			}
		}
	}
}

func getNeighbors(provinces map[string]*prov, mapDataFolder string, mapDataBackupFolder string, color2provID map[string]string) {
	// Open Province File
	provinceFile, err := os.Open(filepath.Join(mapDataFolder, "provinces.png"))
	if err != nil {
		if len(mapDataBackupFolder) > 0 {
			provinceFile, err = os.Open(filepath.Join(mapDataBackupFolder, "provinces.png"))
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println("Did not find provinces.png in \"" + mapDataFolder + "\"")
				fmt.Println("Using backup copy in \"" + mapDataBackupFolder + "\"")
			}
		} else {
			log.Fatal(err)
		}
	}

	provinceMap, err := png.Decode(provinceFile)
	if err != nil {
		log.Fatal(err)
	}

	for _, province := range provinces {
		if province.isValidForHolding {
			for _, pix := range province.pixels {
				r, g, b, _ := provinceMap.At(int(pix.x), int(pix.y)).RGBA()
				unifiedColor := rgbaToUnified(int(r/0x101), int(g/0x101), int(b/0x101))

				neighborPixels := [][]int{{int(pix.x) + 1, int(pix.y) + 1}, {int(pix.x) - 1, int(pix.y) + 1}, {int(pix.x) + 1, int(pix.y) - 1}, {int(pix.x) - 1, int(pix.y) - 1}}
				for _, neighbor := range neighborPixels {
					r, g, b, _ := provinceMap.At(neighbor[0], neighbor[1]).RGBA()
					neighborColor := rgbaToUnified(int(r/0x101), int(g/0x101), int(b/0x101))
					if neighborColor != unifiedColor /*&& provinces[color2provID[neighborColor]].isValidForHolding*/ {
						makeNeighbors(province, provinces[color2provID[neighborColor]])
					}
				}
			}
		}
	}
}

func getAdjacencies(provinces map[string]*prov, mapDataFolder string, mapDataBackupFolder string) {
	// Open Province File
	thisFile, err := os.Open(filepath.Join(mapDataFolder, "adjacencies.csv"))
	if err != nil {
		if len(mapDataBackupFolder) > 0 {
			thisFile, err = os.Open(filepath.Join(mapDataBackupFolder, "adjacencies.csv"))
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println("Did not find adjacencies.csv in \"" + mapDataFolder + "\"")
				fmt.Println("Using backup copy in \"" + mapDataBackupFolder + "\"")
			}
		} else {
			log.Fatal(err)
		}
	}
	// Initialize scanner
	scanner := bufio.NewScanner(thisFile)

	for scanner.Scan() {
		// get next line
		line := scanner.Text()
		fields := strings.Split(line, ";")
		if len(fields) > 2 {
			if fields[2] == "river_large" {
				provinces[fields[0]].isMajorRiverCrossing = true
				provinces[fields[1]].isMajorRiverCrossing = true
			} else if fields[2] == "sea" {
				provinces[fields[0]].isStraitCrossing = true
				provinces[fields[1]].isStraitCrossing = true
			}
			makeNeighbors(provinces[fields[0]], provinces[fields[1]])
		}

	}
}

func makeNeighbors(prov1 *prov, prov2 *prov) {
	if prov1 == nil || prov2 == nil {
		return
	}
	for _, v := range prov1.neighbors {
		if v == prov2.ID {
			return
		}
	}
	prov1.neighbors = append(prov1.neighbors, prov2.ID)
	prov2.neighbors = append(prov2.neighbors, prov1.ID)
}

func getAvgPos(provinces map[string]*prov) {
	for _, prov := range provinces {
		for _, pix := range prov.pixels {
			prov.x += pix.x
			prov.y += pix.y
			prov.elevation += pix.elevation
		}
		prov.x = prov.x / float64(len(prov.pixels))
		prov.y = prov.y / float64(len(prov.pixels))
		prov.elevation = prov.elevation / float64(len(prov.pixels))
	}
}


func getProvinces(mapDataFolder string, mapDataBackupFolder string) (map[string]*prov, map[string]string) {
	// Open Province File
	thisFile, err := os.Open(filepath.Join(mapDataFolder, "definition.csv"))
	if err != nil {
		fmt.Println("Did not find definition.csv in \"" + mapDataFolder + "\"")
		if len(mapDataBackupFolder) > 0 {
			thisFile, err = os.Open(filepath.Join(mapDataBackupFolder, "definition.csv"))
			if err != nil {
				fmt.Println("Did not find definition.csv in \"" + mapDataBackupFolder + "\"")
				log.Fatal(err)
			} else {
				fmt.Println("Using backup copy in \"" + mapDataBackupFolder + "\"")
			}
		} else {
			fmt.Println("No backup map data folder was provided")
			log.Fatal(err)
		}
	}

	// Initialize scanner
	scanner := bufio.NewScanner(thisFile)

	provinces := make(map[string]*prov)
	color2provID := make(map[string]string)

	for scanner.Scan() {
		// get next line
		line := scanner.Text()
		fields := strings.Split(line, ";")
		if len(fields) > 3 {
			color := fields[1] + "." + fields[2] + "." + fields[3]

			var newProv prov
			newProv.ID = fields[0]
			newProv.name = fields[4]
			newProv.color = color
			newProv.pixels = []*pixel{}

			provinces[newProv.ID] = &newProv
			color2provID[color] = newProv.ID
		}
	}
	return provinces, color2provID
}

func getProvinceInfo(mapDataFolder string, mapDataBackupFolder string, provinces map[string]*prov, color2provID map[string]string) map[string]*prov {
	// Open Province File
	provinceFile, err := os.Open(filepath.Join(mapDataFolder, "provinces.png"))
	if err != nil {
		if len(mapDataBackupFolder) > 0 {
			provinceFile, err = os.Open(filepath.Join(mapDataBackupFolder, "provinces.png"))
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println("Did not find provinces.png in \"" + mapDataFolder + "\"")
				fmt.Println("Using backup copy in \"" + mapDataBackupFolder + "\"")
			}
		} else {
			log.Fatal(err)
		}
	}

	provinceMap, err := png.Decode(provinceFile)
	if err != nil {
		log.Fatal(err)
	}

	// Open River File
	riverFile, err := os.Open(filepath.Join(mapDataFolder, "rivers.png"))
	if err != nil {
		if len(mapDataBackupFolder) > 0 {
			riverFile, err = os.Open(filepath.Join(mapDataBackupFolder, "rivers.png"))
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println("Did not find rivers.png in \"" + mapDataFolder + "\"")
				fmt.Println("Using backup copy in \"" + mapDataBackupFolder + "\"")
			}
		} else {
			log.Fatal(err)
		}
	}

	riverMap, err := png.Decode(riverFile)
	if err != nil {
		log.Fatal(err)
	}

	// Open HeightMap File
	heightFile, err := os.Open(filepath.Join(mapDataFolder, "heightmap.png"))
	if err != nil {
		if len(mapDataBackupFolder) > 0 {
			heightFile, err = os.Open(filepath.Join(mapDataBackupFolder, "heightmap.png"))
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println("Did not find heightmap.png in \"" + mapDataFolder + "\"")
				fmt.Println("Using backup copy in \"" + mapDataBackupFolder + "\"")
			}
		} else {
			log.Fatal(err)
		}
	}

	heightMap, err := png.Decode(heightFile)
	if err != nil {
		log.Fatal(err)
	}

	for y := provinceMap.Bounds().Min.Y; y < provinceMap.Bounds().Max.Y; y++ {
		for x := provinceMap.Bounds().Min.X; x < provinceMap.Bounds().Max.X; x++ {
			r, g, b, _ := provinceMap.At(x, y).RGBA()
			unifiedColor := rgbaToUnified(int(r / 0x101), int(g / 0x101) , int(b / 0x101))
			//fmt.Println(unifiedColor)
			provID := color2provID[unifiedColor]


			var pix pixel
			pix.x = float64(x)
			pix.y = float64(y)
			pix.color = unifiedColor
			pix.provID = provID

			// rivermap info
			rr, rg, rb, _ := riverMap.At(x, y).RGBA()
			unifiedColor = rgbaToUnified(int(rr / 0x101), int(rg / 0x101) , int(rb / 0x101))

			if unifiedColor == "255.255.255" {
				pix.class = "land"
			} else if unifiedColor == "255.0.128" {
				pix.class = "sea"
			} else {
				pix.class = "river"
			}

			// heightmap info
			hr, hg, hb, _ := heightMap.At(x, y).RGBA()
			pix.elevation = rgbToElevation(int(hr / 0x101), int(hg / 0x101), int(hb / 0x101))

			if _, ok := provinces[provID]; ok {
				provinces[provID].pixels = append(provinces[provID].pixels, &pix)
			} else {
				fmt.Println("trying to add pixel to province \"" + provID + "\", which does not exist.\nThis error occurs when the pixel color was not found in map_data\\definition.csv.")
			}
		}
	}
	return provinces
}



type pixel struct {
	color string
	provID string
	elevation float64
	class string
	x float64
	y float64
}

func rgbaToUnified(r int, g int, b int) string {
	newr := float64(r)
	newg := float64(g)
	newb := float64(b)
	return strconv.Itoa(int(newr)) + "." + strconv.Itoa(int(newg)) + "." + strconv.Itoa(int(newb))
}

func rgbToElevation(r int, g int, b int) float64 {

	sum := (float64(r) + float64(g) + float64(b)) / float64(3)

	x1 := 19.546
	x2 := 70.043
	y1 := -6.6
	y2 := 1883.0

	slope := (y2 - y1) / (x2 - x1)
	intercept := (-1.0) * slope * x1

	return sum * slope + intercept

}