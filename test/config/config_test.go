package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"imt-atlantique.project.group.fr/meteo-airport/internal/config"
)

type MockConfig struct {
	SomeField string `yaml:"someField"`
}

func (c *MockConfig) Validate() error {
	return nil
}

func TestRetrievePropertiesFromYaml(t *testing.T) {
	tempFile, _ := os.CreateTemp("", "config_test")
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			panic(err)
		}
	}(tempFile.Name())

	_, err := tempFile.WriteString("someField: testValue")
	if err != nil {
		return
	}

	cfg := &MockConfig{}

	err = config.RetrievePropertiesFromYaml(tempFile.Name(), cfg)

	assert.NoError(t, err)
	assert.Equal(t, "testValue", cfg.SomeField)
}

func TestLoadDefaultConfig(t *testing.T) {
	// Create a temporary config file in the expected location
	err := os.MkdirAll("/path/to/config", os.ModePerm)
	if err != nil {
		return
	}
	err = os.WriteFile("/path/to/config/config.yaml", []byte("someField: testValue"), 0644)
	if err != nil {
		return
	}

	cfg := &MockConfig{}

	err = config.LoadDefaultConfig(cfg)

	assert.NoError(t, err)
	assert.Equal(t, "testValue", cfg.SomeField)
}
