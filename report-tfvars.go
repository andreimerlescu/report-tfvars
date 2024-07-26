package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/andreimerlescu/configurable"
)

var SearchPattern *regexp.Regexp

func regex() (*regexp.Regexp, error) {
	var err error
	if SearchPattern == nil {
		SearchPattern, err = regexp.Compile(`variable\s*"(.*?)"\s*\{\s*(?:description\s*=\s*"(.*?)"\s*)?(?:type\s*=\s*([^\n]*?))?\s*(?:default\s*=\s*([\s\S]*?))?\s*\}`)
	}
	return SearchPattern, err
}

func removeFirstCharIfQuote(s string) string {
	if len(s) > 0 && s[0] == '"' {
		return s[1:]
	}
	return s
}

func removeLastCharIfQuote(s string) string {
	if len(s) > 0 && s[len(s)-1] == '"' {
		return s[:len(s)-1]
	}
	return s
}

func clean(in string) string {
	in = strings.ReplaceAll(in, "\n", " ")
	in = strings.ReplaceAll(in, "\t", " ")
	in = strings.Join(strings.Fields(in), " ")

	in = removeFirstCharIfQuote(in)
	in = removeLastCharIfQuote(in)
	return in
}

func processFile(filePath string, regex *regexp.Regexp) {
	tfBytes, readErr := os.ReadFile(filePath)
	if readErr != nil {
		log.Printf("Error reading file %s: %v", filePath, readErr)
		return
	}

	matches := regex.FindAllSubmatch(tfBytes, -1)
	if matches == nil {
		log.Printf("No matches found in file %s.\n", filePath)
		return
	}

	for _, match := range matches {
		if len(match) >= 5 {
			varName := strings.TrimSpace(string(match[1]))
			varDesc := strings.TrimSpace(string(match[2]))
			varType := strings.TrimSpace(string(match[3]))
			varDefault := strings.TrimSpace(string(match[4]))

			defaultParts := strings.Split(varDefault, ` # `)
			if len(defaultParts) >= 2 {
				varDefault = strings.ReplaceAll(defaultParts[0], `"`, ``)
			}

			varName = clean(varName)
			varDesc = clean(varDesc)
			varType = clean(varType)
			varDefault = clean(varDefault)

			if strings.HasPrefix(varDefault, "{") && !strings.HasSuffix(varDefault, "}") {
				varDefault += " }"
			}
			if strings.HasPrefix(varDefault, "[") && !strings.HasSuffix(varDefault, "]") {
				varDefault += " ]"
			}

			fmt.Printf("-var %s='%s'\n", varName, varDefault)
			if varDesc != "" {
				fmt.Printf("        ~> %s value only (%s)\n", varType, varDesc)
			} else {
				fmt.Printf("        ~> %s\n", varType)
			}
		}
	}
}

type Config struct {
	TFFile      *string `json:"file" yaml:"File"`
	TFDirectory *string `json:"directory" yaml:"Directory"`
}

type Application struct {
	cfg    configurable.IConfigurable
	Config Config
}

func main() {
	app := Application{
		cfg: configurable.New(),
	}
	config := Config{
		TFFile:      app.cfg.NewString("file", "", "Path to Terraform file that contains variables."),
		TFDirectory: app.cfg.NewString("dir", "", "Path to Terraform directory to scan for variables."),
	}
	app.Config = config

	var usingConfigFile bool
	configInfo, configInfoErr := os.Stat(filepath.Join(".", "config.yaml"))
	if configInfoErr != nil && configInfo != nil {
		usingConfigFile = false
	} else if configInfoErr != nil && configInfo == nil {
		usingConfigFile = false
	} else if configInfo != nil && configInfo.IsDir() {
		if configInfo.IsDir() {
			log.Println("WARNING: ./config.yaml is a directory not a file and " +
				"the script can auto-load arguments using this file, but cannot " +
				"while it's a directory.")
			usingConfigFile = false
		}
	} else {
		usingConfigFile = true
	}

	if usingConfigFile {
		cfgParseErr := app.cfg.Parse(filepath.Join(".", "config.yaml"))
		if cfgParseErr != nil {
			log.Fatalln(cfgParseErr)
		}
	} else {
		cfgParseErr := app.cfg.Parse("")
		if cfgParseErr != nil {
			log.Fatalln(cfgParseErr)
		}
	}

	// Compile regex
	regex, regexErr := regex()
	if regexErr != nil {
		log.Fatalln(regexErr)
	}

	// Check if the TFFile is provided and process it
	if app.Config.TFFile != nil && len(*app.Config.TFFile) > 0 {
		info, infoErr := os.Stat(*app.Config.TFFile)
		if infoErr != nil {
			log.Fatalln(infoErr)
		}

		if info.IsDir() {
			log.Fatalln("Directory defined for --file when expecting a file.")
		}

		processFile(*app.Config.TFFile, regex)
		return
	}

	// If using --dir instead of --file, scan the directory for .tf files and process each of them individually
	if app.Config.TFDirectory != nil && len(*app.Config.TFDirectory) > 0 {
		err := filepath.Walk(*app.Config.TFDirectory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".tf" {
				processFile(path, regex)
			}
			return nil
		})
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Fatalln("Missing parameter --file or --dir.")
	}
}
