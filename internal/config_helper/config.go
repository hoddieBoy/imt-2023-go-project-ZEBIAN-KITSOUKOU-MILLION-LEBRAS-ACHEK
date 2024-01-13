package config_helper

import (
	"gopkg.in/yaml.v3"
	"imt-atlantique.project.group.fr/meteo-airport/internal/logutil"
	"os"
	"path/filepath"
)

type Config interface {
	Validate() error
}

var defaultConfigFileName = "config.yaml"

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
	return nil
}

func LoadDefaultConfig(cfg Config) error {
	exePath, err := os.Executable()
	if err != nil {
		logutil.Error("Failed to retrieve executable path for loading default config: %v", err)
		return err
	}

	exeDir := filepath.Dir(exePath)
	configPath := filepath.Join(exeDir, defaultConfigFileName)

	return RetrievePropertiesFromYaml(configPath, cfg)
}
