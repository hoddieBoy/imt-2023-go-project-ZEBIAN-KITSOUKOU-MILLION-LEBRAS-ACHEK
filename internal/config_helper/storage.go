package config_helper

import (
	"fmt"
	"imt-atlantique.project.group.fr/meteo-airport/internal/logutil"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt_helper"
)

type StorageConfig struct {
	Storages map[string]map[string]*Storage
	MQTT     *mqtt_helper.MQTTConfig
}

type Storage struct {
	InfluxDB InfluxDBSettings
	CSV      CSVSettings
}

type InfluxDBSettings struct {
	URL          string
	Token        string
	Bucket       string
	Organization string
}

type CSVSettings struct {
	PathDirectory string
	Separator     rune
	TimeFormat    string
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
	logutil.Info("Validating storage config: %v", c.Storages)
	for measurement, storages := range c.Storages {
		storageTypes := make(map[string]bool)

		for _, storage := range storages {
			if err := storage.Validate(); err != nil {
				return err
			}

			storageTypes["influxdb"] = storage.InfluxDB != (InfluxDBSettings{})
			storageTypes["csv"] = storage.CSV != (CSVSettings{})
		}

		if !storageTypes["influxdb"] && !storageTypes["csv"] {
			return fmt.Errorf("no storage type found for measurement %s", measurement)
		}

		if storageTypes["influxdb"] && storageTypes["csv"] {
			return fmt.Errorf("multiple storage types found for measurement %s", measurement)
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

	if c.Separator == 0 {
		return fmt.Errorf("csv separator is empty")
	}

	if c.TimeFormat == "" {
		return fmt.Errorf("csv time format is empty")
	}

	return nil
}

func RetrieveStoragePropertiesFromYaml(filePath string) (*StorageConfig, error) {
	var cfg StorageConfig
	err := RetrievePropertiesFromYaml(filePath, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func LoadDefaultStorageConfig() (*StorageConfig, error) {
	var cfg = StorageConfig{}
	err := LoadDefaultConfig(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
