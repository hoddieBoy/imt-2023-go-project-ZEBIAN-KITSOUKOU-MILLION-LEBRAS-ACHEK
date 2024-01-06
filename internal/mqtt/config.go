package mqtt

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Server struct {
		Port     int    `yaml:"port"`
		Host     string `yaml:"host"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"server"`
}

func RetrievePropertiesFromConfig(filePath string) (*Config, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if decoder.Decode(&cfg) != nil {
		return nil, fmt.Errorf("failed to decode file: %w", err)
	}
	return &cfg, err
}
