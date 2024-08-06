package main

import (
	"fmt"
	"os"
	"slices"
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

func getDependencies(key string, vars map[string]string) []string {
	var res []string
	value := vars[key]
	for k := range vars {
		if strings.Contains(value, fmt.Sprintf("${%s}", k)) || strings.Contains(value, fmt.Sprintf("$%s", k)) {
			res = append(res, k)
		}
	}
	return res
}

func dfs(key string, vars map[string]string, visited map[string]bool, path []string, ordered *[]string) error {
	if visited[key] {
		return nil
	}

	visited[key] = true
	for _, dep := range getDependencies(key, vars) {
		path = append(path, key)
		if slices.Contains(path, dep) {
			return fmt.Errorf("circle detected")
		}
		err := dfs(dep, vars, visited, path, ordered)
		if err != nil {
			return err
		}
	}

	*ordered = append(*ordered, key)
	return nil
}

func resolveVar(value string, envMap map[string]string) string {
	// Replace $$ with a single $ to handle escaping;
	// maybe we don't need this.
	value = strings.ReplaceAll(value, "$$", "$")
	newValue := os.Expand(value, func(key string) string {
		if val, ok := envMap[key]; ok {
			return val
		}
		return os.Getenv(key)
	})
	return newValue
}
func renderEnvironment(envDef EnvironmentDefinition) map[string]string {
	envMap := map[string]string{}

	orderedKeys := []string{}
	visited := map[string]bool{}
	for key := range envDef.Variables {
		path := []string{}
		err := dfs(key, envDef.Variables, visited, path, &orderedKeys)
		if err != nil {
			panic(err)
		}
	}

	for _, key := range orderedKeys {
		envMap[key] = resolveVar(envDef.Variables[key], envMap)
	}

	return envMap
}

func main() {
	for i := 1; i <= 3; i++ {
		envDef, err := readYAMLFile(fmt.Sprintf("environment%d.yaml", i))
		if err != nil {
			fmt.Println("Error reading YAML file:", err)
			return
		}

		renderedEnv := renderEnvironment(envDef)
		for key, value := range renderedEnv {
			fmt.Printf("%s=%s\n", key, value)
		}
		fmt.Println()
	}
}
