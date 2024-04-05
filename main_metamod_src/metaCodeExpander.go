package main

import (
	"fmt"
	"strings"
)

func metaCodeExpander(scriptedEffects map[string]*scriptedEffect, concepts map[string]*concept) map[string]*scriptedEffect {

	for _, scriptedEffect := range scriptedEffects {
		metaCodeExpanderHelper(scriptedEffect, concepts)
	}

	return scriptedEffects
}

func metaCodeExpanderHelper(effect *scriptedEffect, concepts map[string]*concept) {

	var newLines []string
	for _, line := range effect.lines {
		metaCodeElement := containsMetaCodeElement(line)
		if len(metaCodeElement) > 0 {
			editedLines := expandMetaCode(line, metaCodeElement, concepts)
			for _, editedLine := range editedLines {
				newLines = append(newLines, editedLine)
			}
		} else {
			newLines = append(newLines, line)
		}
	}
	effect.lines = newLines

}

func expandMetaCode(line string, metaCodeElement string, concepts map[string]*concept) []string {
	indentLevel := getIndentLevel(line)

	fields := strings.Fields(metaCodeElement)
	conceptName := fields[0]
	target := "&" + conceptName + "&"
	params := make(map[string]string)
	for i := 1; i < len(fields); i++ {
		param := strings.Split(fields[i], ":")
		params[param[0]] = param[1]
	}

	var matchingInstances []*instance
	if _, ok := concepts[conceptName]; ok {
		for _, instance := range concepts[conceptName].instances {
			matching := true
			for paramName, targetVal := range params {
				if instance.params[paramName] != targetVal {
					matching = false
					break
				}
			}
			if matching == true {
				matchingInstances = append(matchingInstances, instance)
			}
		}
	} else {
		fmt.Println("Warning: invalid concept name " + conceptName + " in line: \"" + line + "\"")
	}


	baseLine := getLineWithoutMetaBlock(line)

	var newLines []string
	for _, instance := range matchingInstances {
		newLine := strings.ReplaceAll(baseLine, target, instance.name)
		newLine = applyIndent(newLine, indentLevel)
		newLines = append(newLines, newLine)
	}

	return newLines
}

func getLineWithoutMetaBlock(line string) string {
	runes := []rune(line)
	target := []rune("^")[0]
	highestIndex := 0
	for i, rune := range runes {
		if rune == target {
			highestIndex = i
		}
	}
	runes = runes[highestIndex+1:]

	return string(runes)
}

func containsMetaCodeElement(line string) string {
	if strings.Count(line, "^") < 1 {
		return ""
	} else {
		runes := []rune(line)
		metaCodeStart := 0
		metaCodeEnd := 0
		counter := 0
		target := []rune("^")[0]
		for i, rune := range runes {
			if rune == target {
				counter++
				if counter == 1 {
					metaCodeStart = i+1
				} else if counter == 2 {
					metaCodeEnd = i
				}
			}
		}
		return line[metaCodeStart:metaCodeEnd]
	}

}