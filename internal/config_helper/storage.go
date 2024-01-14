package config_helper

import (
	"fmt"

	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
)

type StorageConfig struct {
	Storages map[string]Storage
	MQTT     mqtt.Config
}

type Storage struct {
	InfluxDB InfluxDBSettings `yaml:"influxdb"`
	CSV      CSVSettings      `yaml:"csv"`
}

type InfluxDBSettings struct {
	URL          string `yaml:"url"`
	Token        string `yaml:"token"`
	Bucket       string `yaml:"bucket"`
	Organization string `yaml:"organization"`
}

type CSVSettings struct {
	PathDirectory string `yaml:"path_directory"`
	Separator     string `yaml:"separator"`
	TimeFormat    string `yaml:"time_format"`
}

func (c *Storage) Validate() error {

	if c.InfluxDB == (InfluxDBSettings{}) && c.CSV == (CSVSettings{}) {
		return fmt.Errorf("influxdb and csv settings are empty")
	}

	if c.InfluxDB != (InfluxDBSettings{}) {
		if err := c.InfluxDB.Validate(); err != nil {
			return err
		}
	}

	if c.CSV != (CSVSettings{}) {
		if err := c.CSV.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (c *StorageConfig) Validate() error {
	for _, storages := range c.Storages {
		// For each measurement, check if the storage is valid
		if err := storages.Validate(); err != nil {
			return err
		}
	}

	if err := c.MQTT.Validate(); err != nil {
		return err
	}

	return nil
}

func (c *InfluxDBSettings) Validate() error {
	if c.URL == "" {
		return fmt.Errorf("influxdb url is empty")
	}

	if c.Token == "" {
		return fmt.Errorf("influxdb token is empty")
	}

	if c.Bucket == "" {
		return fmt.Errorf("influxdb bucket is empty")
	}

	if c.Organization == "" {
		return fmt.Errorf("influxdb organization is empty")
	}

	return nil
}

func (c *CSVSettings) Validate() error {
	if c.PathDirectory == "" {
		return fmt.Errorf("csv path directory is empty")
	}

	if c.Separator == "" && len(c.Separator) > 1 {
		return fmt.Errorf("csv separator is empty or more than one character")
	}

	if c.TimeFormat == "" {
		return fmt.Errorf("csv time format is empty")
	}

	return nil
}

func LoadDefaultStorageConfig() (*StorageConfig, error) {
	cfg := &StorageConfig{}
	err := LoadDefaultConfig(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
