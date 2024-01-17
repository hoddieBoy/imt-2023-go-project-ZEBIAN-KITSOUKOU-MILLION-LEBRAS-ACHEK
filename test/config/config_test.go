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

func createTempFileWithContent(content string) (string, error) {
	tempFile, err := os.CreateTemp("", "config_test")
	if err != nil {
		return "", err
	}

	_, err = tempFile.WriteString(content)
	if err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}

func TestRetrievePropertiesFromYaml(t *testing.T) {
	tempFileName, err := createTempFileWithContent("someField: testValue")
	if err != nil {
		t.Fatal(err)
	}

	defer func(name string) {
		err = os.Remove(name)
		if err != nil {
			t.Fatal(err)
		}
	}(tempFileName)

	cfg := &MockConfig{}

	err = config.RetrievePropertiesFromYaml(tempFileName, cfg)

	assert.NoError(t, err)
	assert.Equal(t, "testValue", cfg.SomeField)
}

func TestLoadDefaultConfig(t *testing.T) {
	// Create a temporary config file in the expected location
	err := os.MkdirAll("/path/to/config", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	tempFileName, err := createTempFileWithContent("someField: testValue")
	if err != nil {
		t.Fatal(err)
	}

	defer func(name string) {
		err = os.Remove(name)

		if err != nil {
			t.Fatal(err)
		}
	}(tempFileName)

	cfg := &MockConfig{}

	err = config.LoadDefaultConfig(cfg)

	assert.NoError(t, err)
	assert.Equal(t, "testValue", cfg.SomeField)
}
