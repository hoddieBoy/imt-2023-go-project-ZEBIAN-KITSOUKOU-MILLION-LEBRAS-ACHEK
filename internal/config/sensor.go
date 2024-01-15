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
		Unit       string `yaml:"unit"`
		TimeFormat string `yaml:"time_format"`
	} `yaml:"sensor"`
	Broker struct {
		Client       *mqtt.Config `yaml:"client"`
		PublishTopic string       `yaml:"publish_topic"`
	} `yaml:"broker"`
}

func (c *SensorConfig) Validate() error {
	if err := c.Broker.Client.Validate(); err != nil {
		return err
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

	return nil
}

func LoadDefaultSensorConfig() (*SensorConfig, error) {
	var cfg SensorConfig
	err := LoadDefaultConfig(&cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
