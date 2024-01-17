package config

import (
	"fmt"

	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
)

type Storage struct {
	Settings map[string]Setting `yaml:"settings"`
	Broker   struct {
		Config   mqtt.Config `yaml:"config"`
		ClientID string      `yaml:"client_id"`
	} `yaml:"broker"`
}

type Setting struct {
	InfluxDB InfluxDBSettings `yaml:"influxdb"`
	CSV      CSVSettings      `yaml:"csv"`
	Qos      byte
	Topic    string
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

func (s *Setting) Validate() error {
	if s.InfluxDB == (InfluxDBSettings{}) && s.CSV == (CSVSettings{}) {
		return fmt.Errorf("influxdb and csv settings are empty")
	}

	if s.Qos > 2 {
		return fmt.Errorf("qos must be between 0 and 2")
	}

	if s.InfluxDB != (InfluxDBSettings{}) {
		if err := s.InfluxDB.Validate(); err != nil {
			return err
		}
	}

	if s.CSV != (CSVSettings{}) {
		if err := s.CSV.Validate(); err != nil {
			return err
		}
	}

	if s.Topic == "" {
		return fmt.Errorf("topic is empty")
	}

	return nil
}

func (c *Storage) Validate() error {
	if len(c.Settings) == 0 {
		return fmt.Errorf("settings is empty")
	}

	for _, settings := range c.Settings {
		if err := settings.Validate(); err != nil {
			return err
		}
	}

	if c.Broker.ClientID == "" {
		return fmt.Errorf("client id is empty")
	}

	return c.Broker.Config.Validate()
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
