package config

import (
	"fmt"

	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
)

type SensorConfig struct {
	Setting struct {
		AirportID  string `yaml:"airport_id"`
		SensorID   int64  `yaml:"sensor_id"`
		Topic      string `yaml:"topic"`
		Type       string `yaml:"type"`
		Unit       string `yaml:"unit"`
		TimeFormat string `yaml:"time_format"`
	} `yaml:"sensor"`
	Broker struct {
		Client mqtt.Config `yaml:"client"`
		Qos    byte        `yaml:"qos"`
	} `yaml:"broker"`
}

func (c *SensorConfig) Validate() error {
	if err := c.Broker.Client.Validate(); err != nil {
		return err
	}

	if c.Broker.Qos < 0 || c.Broker.Qos > 2 {
		return fmt.Errorf("qos must be between 0 and 2")
	}

	if c.Setting.AirportID == "" {
		return fmt.Errorf("airport id is empty")
	}

	if c.Setting.SensorID == 0 {
		return fmt.Errorf("sensor id is empty")
	}

	if c.Setting.Topic == "" {
		return fmt.Errorf("topic is empty")
	}

	if c.Setting.Unit == "" {
		return fmt.Errorf("unit is empty")
	}

	if c.Setting.TimeFormat == "" {
		return fmt.Errorf("time format is empty")
	}

	if c.Setting.Type == "" {
		return fmt.Errorf("type is empty")
	}

	return nil
}

func LoadDefaultSensorConfig() (*SensorConfig, error) {
	var cfg SensorConfig
	err := LoadDefaultConfig(&cfg)

	return &cfg, err
}
