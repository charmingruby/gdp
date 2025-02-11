package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	CONFIG_FILE = "config.yml"
)

type Config struct {
	Client ClientConfig `yaml:"client"`
	Server ServerConfig `yaml:"server"`
}

type ClientConfig struct {
	ServerPort int `yaml:"server-port"`
}

type ServerConfig struct {
	Port      int             `yaml:"port"`
	Threshold ThresholdConfig `yaml:"thresholds"`
}

type ThresholdConfig struct {
	PackageLoss float32 `yaml:"package-loss"`
}

func New() (*Config, error) {
	configFileExists := validateExistence()
	if !configFileExists {
		return nil, fmt.Errorf("config file does not exist")
	}

	return parseConfigFile()
}

func validateExistence() bool {
	_, err := os.Stat(CONFIG_FILE)

	isConfigFileInvalid := os.IsNotExist(err)

	return !isConfigFileInvalid
}

func parseConfigFile() (*Config, error) {
	yamlData, err := os.ReadFile(CONFIG_FILE)
	if err != nil {
		return &Config{}, err
	}

	var config Config
	if err := yaml.Unmarshal(yamlData, &config); err != nil {
		return &Config{}, err
	}

	return &config, nil
}
