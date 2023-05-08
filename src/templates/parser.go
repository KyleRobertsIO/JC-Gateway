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

type YamlConfig struct {
	Image                string        `yaml:"image"`
	Resources            YamlResources `yaml:"resources"`
	Ports                []string      `yaml:"ports"`
	Commands             []string      `yaml:"commands"`
	EnvironmentVariables []string      `yaml:"environment_variables"`
}

func Parse(templateName string) (*YamlConfig, error) {
	// Open the targeted job template file
	b, readErr := os.ReadFile("./job_templates/basic.yml")
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
