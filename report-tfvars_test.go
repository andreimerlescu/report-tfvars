package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/andreimerlescu/configurable"
)

func TestRegex(t *testing.T) {
	r, err := regex()
	if err != nil {
		t.Fatalf("Failed to compile regex: %v", err)
	}
	if r == nil {
		t.Fatalf("Expected a valid regex pattern, got nil")
	}
}

func TestRemoveFirstCharIfQuote(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{`"hello`, `hello`},
		{`hello`, `hello`},
		{`"`, ``},
	}

	for _, test := range tests {
		result := removeFirstCharIfQuote(test.input)
		if result != test.output {
			t.Errorf("Expected %v, got %v", test.output, result)
		}
	}
}

func TestRemoveLastCharIfQuote(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{`hello"`, `hello`},
		{`hello`, `hello`},
		{`"`, ``},
	}

	for _, test := range tests {
		result := removeLastCharIfQuote(test.input)
		if result != test.output {
			t.Errorf("Expected %v, got %v", test.output, result)
		}
	}
}

func TestClean(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"  hello   world  ", "hello world"},
		{"\thello\tworld\t", "hello world"},
		{"\nhello\nworld\n", "hello world"},
		{`"hello"`, "hello"},
	}

	for _, test := range tests {
		result := clean(test.input)
		if result != test.output {
			t.Errorf("Expected %v, got %v", test.output, result)
		}
	}
}

func TestProcessFile(t *testing.T) {
	testFile := "vars.tf"
	// Create the test file with some content
	content := `
variable "color" {
	description = "Deployment group color or label or tag"
	type = string
	default = "blue"
}

variable "complex" {
	description = "this is a complex description"
	type = map(string)
	default = {
		"alfred" = "Alfred Wilson"
	}
}

variable "public_key" {
	description = "Value of the SSH public key"
	type = string
	default = "value"
}
`

	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	defer os.Remove(testFile) // Clean up the test file after running the test

	r, err := regex()
	if err != nil {
		t.Fatalf("Failed to compile regex: %v", err)
	}

	processFile(testFile, r)
	// Normally, you would capture the output and assert it against expected values
	// For simplicity, we are just running the function to ensure it executes without error
}

func TestMain(t *testing.T) {
	// Set up the environment variables
	os.Setenv("CI_PROJECT_DIR", ".")
	defer os.Unsetenv("CI_PROJECT_DIR")

	// Create a temporary directory and change to it
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create a dummy vars.tf file
	if err := os.WriteFile("vars.tf", []byte(`variable "test" { type = string default = "value" }`), 0644); err != nil {
		t.Fatalf("Failed to create vars.tf file: %v", err)
	}

	// Set up the config file
	f := filepath.Join(dir, "vars.tf")
	config := Config{
		TFFile:      &f,
		TFDirectory: &dir,
	}

	app := Application{
		cfg:    configurable.New(),
		Config: config,
	}

	// Test the main function logic
	if len(*app.Config.TFFile) == 0 {
		t.Fatalf("Invalid value for --file.")
	}

	info, infoErr := os.Stat(*app.Config.TFFile)
	if infoErr != nil {
		t.Fatalf("Failed to stat --file: %v", infoErr)
	}

	if info.IsDir() {
		t.Fatalf("Directory defined for --file when expecting a file.")
	}

	r, err := regex()
	if err != nil {
		t.Fatalf("Failed to compile regex: %v", err)
	}

	processFile(*app.Config.TFFile, r)
}

func TestConfigurable(t *testing.T) {
	app := Application{
		cfg: configurable.New(),
	}
	config := Config{
		TFFile:      app.cfg.NewString("file", filepath.Join(".", "vars.tf"), "Path to Terraform file that contains variables."),
		TFDirectory: app.cfg.NewString("dir", filepath.Join(".", "terraform"), "Path to Terraform directory to scan for variables."),
	}
	app.Config = config

	if app.Config.TFFile == nil {
		t.Fatal("TFFile should not be nil")
	}

	if *app.Config.TFFile != filepath.Join(".", "vars.tf") {
		t.Fatalf("Expected %s, got %s", filepath.Join(".", "vars.tf"), *app.Config.TFFile)
	}
}
