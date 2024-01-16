package config

import (
	"fmt"

	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
)

type Storage struct {
	Settings map[string]struct {
		InfluxDB InfluxDBSettings `yaml:"influxdb"`
		CSV      CSVSettings      `yaml:"csv"`
		Qos      byte
		Topic    string
	} `yaml:"settings"`
	Broker struct {
		Config   mqtt.Config `yaml:"config"`
		ClientID string      `yaml:"client_id"`
	} `yaml:"broker"`
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
	if len(c.Settings) == 0 {
		return fmt.Errorf("settings is empty")
	}

	for _, settings := range c.Settings {
		if settings.InfluxDB == (InfluxDBSettings{}) && settings.CSV == (CSVSettings{}) {
			return fmt.Errorf("influxdb and csv settings are empty")
		}

		if settings.Qos < 0 || settings.Qos > 2 {
			return fmt.Errorf("qos must be between 0 and 2")
		}

		if settings.InfluxDB != (InfluxDBSettings{}) {
			if err := settings.InfluxDB.Validate(); err != nil {
				return err
			}
		}

		if settings.CSV != (CSVSettings{}) {
			if err := settings.CSV.Validate(); err != nil {
				return err
			}
		}

		if settings.Topic == "" {
			return fmt.Errorf("topic is empty")
		}
	}

	if err := c.Broker.Config.Validate(); err != nil {
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

func LoadDefaultStorageConfig() (*Storage, error) {
	cfg := &Storage{}
	err := LoadDefaultConfig(cfg)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}
