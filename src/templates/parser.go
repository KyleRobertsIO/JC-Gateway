package templates

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type YamlResources struct {
	CPU    int `yaml:"cpu"`
	Memory int `yaml:"memory"`
}

type YamlEnvironmentVariable struct {
	Secure bool   `json:"secure"`
	Name   string `json:"name"`
	Value  string `json:"value"`
}

type YamlConfig struct {
	ContainerName        string                    `yaml:"container_name"`
	OperatingSystem      string                    `yaml:"os"`
	Image                string                    `yaml:"image"`
	Resources            YamlResources             `yaml:"resources"`
	Ports                []string                  `yaml:"ports"`
	Commands             []string                  `yaml:"commands"`
	EnvironmentVariables []YamlEnvironmentVariable `yaml:"environment_variables"`
}

func Parse(templateName string) (*YamlConfig, error) {
	// Open the targeted job template file
	templatePath := fmt.Sprintf("./job_templates/%s.yml", templateName)
	b, readErr := os.ReadFile(templatePath)
	if readErr != nil {
		return nil, fmt.Errorf(
			"failed to read job template file; %s",
			readErr.Error(),
		)
	}
	// Parse the YAML config from targeted template file
	yamlStr := string(b)
	config := YamlConfig{}
	yamlParseErr := yaml.Unmarshal([]byte(yamlStr), &config)
	if yamlParseErr != nil {
		return nil, fmt.Errorf(
			"failed to parse job template yaml; %s",
			yamlParseErr.Error(),
		)
	}
	return &config, nil
}
