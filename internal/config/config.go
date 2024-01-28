package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config interface {
	Validate() error
}

var defaultConfigFileName = "config/config.yaml"

func SetDefaultConfigFileName(name string) {
	defaultConfigFileName = name
}

func RetrievePropertiesFromYaml(filePath string, cfg Config) error {
	file, err := os.ReadFile(filePath)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, cfg)

	if err != nil {
		return err
	}

	return cfg.Validate()
}

func LoadDefaultConfig(cfg Config) error {
	exePath, err := os.Executable()
	// _, filename, _, ok := runtime.Caller(0)
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	exeDir := filepath.Dir(exePath)
	configPath := filepath.Join(exeDir, defaultConfigFileName)

	return RetrievePropertiesFromYaml(configPath, cfg)
}
