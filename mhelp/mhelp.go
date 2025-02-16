package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type helpConfig struct {
	filepath string
	sort     bool
}

type variable struct {
	name          string
	value         string
	documentation string
}

type target struct {
	name          string
	documentation string
}

const colourReset = "\033[0m"     // reset
const colourTarget = "\033[36m"   // cyan
const colourVariable = "\033[93m" // yellow
const colourDefault = "\033[90m"  // grey

func isVariable(tokens []string) bool {
	return strings.Contains(tokens[0], "=") ||
		strings.Contains(tokens[0], "?=") ||
		strings.Contains(tokens[0], ":=") ||
		strings.Contains(tokens[0], "+=")
}

func cleanAndBuildTarget(tokens []string) target {
	return target{
		name:          strings.SplitN(tokens[0], ":", 2)[0],
		documentation: strings.TrimSpace(tokens[1]),
	}
}

func cleanAndBuildVariable(tokens []string) variable {
	tokensAssignment := strings.SplitN(tokens[0], "=", 2)
	return variable{
		name:          strings.TrimSpace(strings.TrimRight(tokensAssignment[0], ":?!+")),
		value:         strings.TrimSpace(tokensAssignment[1]),
		documentation: strings.TrimSpace(tokens[1]),
	}
}

func main() {
	var hc helpConfig
	flag.StringVar(&hc.filepath, "filepath", "./Makefile", "The filepath of your Makefile.")
	flag.BoolVar(&hc.sort, "sort", true, "Sort variable and target names alphabetically.")
	flag.Parse()

	file, err := os.Open(hc.filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var line string
	var targets []target
	var variables []variable
	var tokens []string

	// collect and parse Makefile lines
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
		tokens = strings.SplitN(line, "##", 2)
		if len(tokens) < 2 || len(tokens[0]) < 1 {
			continue
		} else if isVariable(tokens) {
			variables = append(variables, cleanAndBuildVariable(tokens))
		} else {
			targets = append(targets, cleanAndBuildTarget(tokens))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// sort variables and targets
	if hc.sort {
		sort.Slice(targets, func(i, j int) bool {
			return targets[i].name < targets[j].name
		})
		sort.Slice(variables, func(i, j int) bool {
			return variables[i].name < variables[j].name
		})
	}

	// build the output to our liking
	output := "Usage:\n" +
		"  make" +
		string(colourTarget) + " [target]" + string(colourReset) +
		string(colourVariable) + " [variables]" + string(colourReset) + "\n"
	output += "\nTarget(s):\n"
	for _, t := range targets {
		output += fmt.Sprintf(
			"  %-30s %s",
			string(colourTarget)+t.name+string(colourReset), t.documentation+"\n",
		)
	}
	output += "\nVariable(s):\n"
	for _, v := range variables {
		output += fmt.Sprintf(
			"  %-30s %s %s",
			string(colourVariable)+v.name+string(colourReset), v.documentation,
			string(colourDefault)+"(default: "+v.value+")"+string(colourReset)+"\n",
		)
	}
	output += "\nExample(s):\n" +
		"  make\n" +
		"  make help\n"

	fmt.Println(output)
}
