package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type EnvironmentDefinition struct {
	Variables map[string]string `yaml:"environment"`
}

func readYAMLFile(filePath string) (EnvironmentDefinition, error) {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return EnvironmentDefinition{}, err
	}

	var env EnvironmentDefinition
	err = yaml.Unmarshal(yamlFile, &env)
	if err != nil {
		return EnvironmentDefinition{}, err
	}

	return env, nil
}

func renderEnvironment(envDef EnvironmentDefinition) map[string]string {
	envMap := map[string]string{}
	unresolved := map[string]string{}

	// Copy the environment variables to unresolved.
	for key, value := range envDef.Variables {
		unresolved[key] = value
	}

	// Function to resolve a single variable value.
	resolveVar := func(value string) string {
		// Replace $$ with a single $ to handle escaping.
		// Maybe we don't need this.
		value = strings.ReplaceAll(value, "$$", "$")

		newValue := os.Expand(value, func(key string) string {
			if val, ok := envMap[key]; ok {
				return val
			}
			return os.Getenv(key)
		})
		return newValue
	}

	for len(unresolved) > 0 {
		// Track if any variables are resolved in this iteration.
		resolvedInThisIteration := false

		for key, value := range unresolved {
			newValue := resolveVar(value)

			// If no change was made, it means all dependencies are resolved.
			if newValue == value {
				envMap[key] = value
				delete(unresolved, key)
				resolvedInThisIteration = true
			} else {
				unresolved[key] = newValue
			}
		}

		// If no variables were resolved in this iteration, it indicates a circular dependency.
		if !resolvedInThisIteration {
			fmt.Println("Circular dependency detected. Unable to resolve further.")
			break
		}
	}

	return envMap
}

func main() {
	envDef, err := readYAMLFile("environment.yaml")
	if err != nil {
		fmt.Println("Error reading YAML file:", err)
		return
	}

	renderedEnv := renderEnvironment(envDef)
	for key, value := range renderedEnv {
		fmt.Printf("%s=%s\n", key, value)
	}
}
