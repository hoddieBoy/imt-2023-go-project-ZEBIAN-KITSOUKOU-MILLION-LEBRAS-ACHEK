package config_helper

import "fmt"

type ApiConfig struct {
	Api struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"api"`

	Database InfluxDBSettings `yaml:"database"`
}

func (c *ApiConfig) Validate() error {
	if c.Api.Host == "" {
		return fmt.Errorf("host is empty")
	}

	if c.Api.Port == "" {
		return fmt.Errorf("port is empty")
	}

	if err := c.Database.Validate(); err != nil {
		return err
	}

	return nil
}

func RetrieveApiPropertiesFromYaml(filePath string) (*ApiConfig, error) {
	var cfg ApiConfig
	err := RetrievePropertiesFromYaml(filePath, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func LoadDefaultApiConfig() (*ApiConfig, error) {
	var cfg ApiConfig
	err := LoadDefaultConfig(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
