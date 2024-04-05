package main

import (
	"strings"
)

func scriptedEffectsExpander(scriptedEffects map[string]*scriptedEffect, concepts map[string]*concept) map[string]*scriptedEffect {

		for _, scriptedEffect := range scriptedEffects {
			expandScriptedEffect(scriptedEffect, scriptedEffects, concepts)
		}

	return scriptedEffects
}

func expandScriptedEffect(effect *scriptedEffect, allEffects map[string]*scriptedEffect, concepts map[string]*concept) {
	foundUnexpandedScriptedEffect := true
	iterations := 0
	for foundUnexpandedScriptedEffect == true && iterations < 100 {
		foundUnexpandedScriptedEffect = false

		var newLines []string
		for _, line := range effect.lines {
			containsEffect, effectName := containsScriptedEffect(line, allEffects)
			if containsEffect {
				editedLines := expandLine(line, allEffects[effectName], concepts)
				for _, editedLine := range editedLines {
					newLines = append(newLines, editedLine)
				}
				foundUnexpandedScriptedEffect = true
			} else {
				newLines = append(newLines, line)
			}
		}
		effect.lines = newLines
		iterations++
	}

}

func expandLine(line string, effect *scriptedEffect, concepts map[string]*concept) []string {
	indentLevel := getIndentLevel(line)
	var editedLines []string

	args := getArgs(line)

	for _, nextLine := range effect.lines {
		editedLine := nextLine
		for k, v := range args {

			editedLine = strings.Replace(editedLine, k, v, -1)
		}
		editedLines = append(editedLines, editedLine)
	}

	for i, nextLine := range editedLines {
		lineEdited := false

		var matches []string
		runes := []rune(nextLine)
		target := []rune("*")[0]
		counter := 0
		start := 0
		end := 0
		for i, rune := range runes {
			if rune == target {
				counter++
				if counter % 2 == 1 {
					start = i
				} else {
					end = i+1
					matches = append(matches, string(runes[start:end]))
				}
			}
		}
		for _, match := range matches {
			property := strings.ReplaceAll(match, "*", "")
			splitProperty := strings.Split(property, ".")
			if len(splitProperty) > 1 {
				instanceName := splitProperty[0]
				instanceField := splitProperty[1]
				for _, concept := range concepts {
					for _, instance := range concept.instances {
						if instance.name == instanceName {
							replacement := instance.params[instanceField]
							nextLine = strings.ReplaceAll(nextLine, match, replacement)
							lineEdited = true
						}
					}
				}
			}
		}
		if lineEdited {
			editedLines[i] = nextLine
		}
	}

	for i, nextLine := range editedLines {
		editedLines[i] = applyIndent(nextLine, indentLevel)
	}

	return editedLines
}

func fields2line(fields []string) string {
	line := fields[0]
	for i := 1; i < len(fields); i++ {
		line += " " + fields[i]
	}
	return line
}

func getArgs(line string) map[string]string {
	args := make(map[string]string)
	fields := strings.Fields(line)

	for i := 1; i < len(fields)-1; i++ {
		if fields[i] == "=" && fields[i+1] != "{" {
			args["$" + fields[i-1] + "$"] = fields[i+1]
		}
	}
	return args
}

func containsScriptedEffect(line string, scriptedEffects map[string]*scriptedEffect) (bool, string) {
	fields := strings.Fields(line)
	// check if line is a scripted effect
	for _, scriptedEffect := range scriptedEffects {
		if fields[0] == scriptedEffect.name { return true, scriptedEffect.name }
	}
	return false, ""
}

func applyIndent(line string, targetIndentLevel int) string {

	indentsToApply := targetIndentLevel - 4

	for i := 0; i < indentsToApply; i++ {
		line = " " + line
	}

	return line
}

func getIndentLevel(line string) int {
	i := 0
	for _, runeValue := range line {
		if runeValue == ' ' {
			i++
		} else if runeValue == []rune("\t")[0] {
			i+=4
		} else {
			break
		}
	}
	return i
}