package webhook

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	inner map[string]*ProjectDetails
}

type ProjectDetails struct {
	Branch  string `yaml:"branch"`
	Command string `yaml:"command"`
}

func (c *Config) Print() {
	for name, details := range c.inner {
		fmt.Printf("%s:\n", name)
		fmt.Printf("  branch -> %s\n", details.Branch)
		fmt.Printf("  command -> \"%s\"\n", details.Command)
	}
}

func (c *Config) get(projectName string) (*ProjectDetails, error) {
	project, ok := c.inner[projectName]
	if !ok {
		return nil, errors.New("project doesn't exist: " + projectName)
	}
	return project, nil
}

type parsedConfig struct {
	Projects []struct {
		Repository string `yaml:"repository"`
		Branch     string `yaml:"branch"`
		Command    string `yaml:"command"`
	} `yaml:"projects"`
}

func ParseConfig(config string) *Config {
	yml, err := os.ReadFile(config)
	if err != nil {
		panic("error reading config: " + err.Error())
	}

	var parsed parsedConfig

	if err := yaml.Unmarshal(yml, &parsed); err != nil {
		panic("error parsing config: " + err.Error())
	}

	result := &Config{inner: make(map[string]*ProjectDetails)}

	for _, p := range parsed.Projects {
		result.inner[p.Repository] = &ProjectDetails{
			Branch:  p.Branch,
			Command: p.Command,
		}
	}

	return result
}
