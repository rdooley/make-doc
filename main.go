package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/willf/pad"
	"gopkg.in/alecthomas/kingpin.v2"
)

const version = "0.0.1"
const targetRegex = `^((?P<name>(-|_|\w)+)|(\${(?P<varname>(-|_|\w)+)}))\s*:.*\#\#+\s*@(?P<category>(\w+))\s+(?P<doc>(.*))`
const varRegex = `^(?P<name>(\w)+)\s*(\?)?=\s*(?P<default>([^\#]+))(\s+\#\#+\s*@(?P<category>(\w+))(:)?\s+(?P<doc>(.*))|\s*$)`

var (
	parseVars = kingpin.Flag("variables", "optional flag to parse variables").Bool()
	makefiles = kingpin.Arg("makefiles", "path to makefiles").Required().ExistingFiles()
	_         = kingpin.Version(version)

	variables = VarInfo{}
	varRe     = regexp.MustCompile(varRegex)
	targets   = MakeTargets{}
	vargets   = MakeTargets{}
	targetRe  = regexp.MustCompile(targetRegex)
)

type VarInfo map[string]string

type MakeTarget struct {
	name  string
	cat   string
	doc   string
	def   string
	isVar bool
}

type MakeTargets map[string][]MakeTarget

func parseLine(line string, re *regexp.Regexp) (target MakeTarget) {
	if !re.MatchString(line) {
		return target
	}
	var cat string
	var def string
	var doc string
	var name string
	var isVar bool
	matches := re.FindStringSubmatch(line)
	for i, n := range re.SubexpNames() {
		switch n {
		case "name":
			if name == "" {
				name = matches[i]
				isVar = false
			}
		case "varname":
			if name == "" {
				name = matches[i]
				isVar = true
			}
		case "def":
			def = strings.TrimSpace(matches[i])
		case "doc":
			doc = strings.TrimSpace(matches[i])
		case "category":
			cat = matches[i]
		}

	}
	if name != "" {
		_, ok := variables[name]
		if !ok {
			variables[name] = def
		}
	}
	return MakeTarget{
		name:  name,
		cat:   cat,
		doc:   doc,
		def:   def,
		isVar: isVar,
	}
}

func parseTargets(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		target := parseLine(line, targetRe)
		if target.cat != "" {
			targets[target.cat] = append(targets[target.cat], target)
		}
		varget := parseLine(line, varRe)
		if varget.cat != "" {
			vargets[varget.cat] = append(vargets[varget.cat], varget)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func printHelp(cats MakeTargets) {
	var colSize int
	for _, rules := range cats {
		for _, rule := range rules {
			l := len(rule.name)
			if l > colSize {
				colSize = l
			}
		}
	}
	for cat, rules := range cats {
		fmt.Printf("\n%s:\n", cat)
		for _, rule := range rules {
			defVal := rule.def
			if defVal != "" {
				defVal = " Default : " + defVal
			}
			fmt.Printf("\t%s\t%s.%s\n", pad.Right(rule.name, colSize, " "), rule.doc, defVal)
		}
	}
}

func main() {
	kingpin.Parse()
	for _, m := range *makefiles {
		parseTargets(m)
	}
	// subVariables()
	if *parseVars {
		printHelp(vargets)
	} else {
		fmt.Println("Usage: make [target] [VARIABLE=value]")
		printHelp(targets)
	}

}
