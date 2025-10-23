package io

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// GetMode returns the execution mode from environment variable
func GetMode() string {
	return os.Getenv("MODE")
}

// GetEnvironment returns the execution environment configuration from environment variable
func GetEnvironment() map[string]any {
	envJSON := os.Getenv("ENVIRONMENT")
	if envJSON == "" {
		return make(map[string]any)
	}

	var env map[string]any
	if err := json.Unmarshal([]byte(envJSON), &env); err != nil {
		fmt.Printf("Error parsing environment JSON: %v\n", err)
		return make(map[string]any)
	}

	return env
}

// GetInputs parses the --inputs YAML argument or --inputs-file and returns a map[string]any
func GetInputs() map[string]any {
	var inputsYAML string
	var inputsFile string
	flag.StringVar(&inputsYAML, "inputs", "", "YAML string containing all inputs")
	flag.StringVar(&inputsFile, "inputs-file", "", "Path to YAML file containing all inputs")
	flag.Parse()

	inputs := make(map[string]any)

	// Priority: file input over inline input (for large data)
	if inputsFile != "" {
		// Read inputs from file
		yamlBytes, err := os.ReadFile(inputsFile)
		if err != nil {
			fmt.Printf("Error reading inputs file %s: %v\n", inputsFile, err)
			return make(map[string]any)
		}
		if err := yaml.Unmarshal(yamlBytes, &inputs); err != nil {
			fmt.Printf("Error parsing YAML from file %s: %v\n", inputsFile, err)
			return make(map[string]any)
		}
	} else if inputsYAML != "" {
		// Read inputs from inline YAML
		if err := yaml.Unmarshal([]byte(inputsYAML), &inputs); err != nil {
			fmt.Printf("Error parsing inline YAML: %v\n", err)
			return make(map[string]any)
		}
	}

	return inputs
}

// SetOutputs converts a map to YAML format and prints it as outputs
func SetOutputs(outputs map[string]any) {
	// Create a wrapper map with "outputs" key
	wrapper := map[string]any{
		"outputs": outputs,
	}

	// Marshal to YAML
	yamlBytes, err := yaml.Marshal(wrapper)
	if err != nil {
		// Fallback to simple format if marshaling fails
		fmt.Println("outputs:")
		for key, value := range outputs {
			fmt.Printf("  %s: %v\n", key, value)
		}
		return
	}

	// Print the YAML output
	fmt.Print(string(yamlBytes))
}
