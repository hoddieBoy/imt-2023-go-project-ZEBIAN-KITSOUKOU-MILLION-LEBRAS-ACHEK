package config_helper

import (
	"fmt"

	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
)

type AlertConfig struct {
	Broker       mqtt.Config            `yaml:"broker"`
	SensorsAlert map[string]SensorAlert `yaml:"sensors_alert"`
}

type SensorAlert struct {
	IncomingTopic string `yaml:"incoming_topic"`
	OutgoingTopic string `yaml:"outgoing_topic"`
	LowerBound    int    `yaml:"lower_bound"`
	HigherBound   int    `yaml:"higher_bound"`
}

func (c *AlertConfig) Validate() error {
	// Check if the MQTTConfig is valid
	if err := c.Broker.Validate(); err != nil {
		return err
	}

	// Check if the incoming topic is valid
	for _, sensor := range c.SensorsAlert {
		// Check if the sensor type is valid
		if sensor.IncomingTopic == "" {
			return fmt.Errorf("incoming topic is empty")
		}

		if sensor.OutgoingTopic == "" {
			return fmt.Errorf("outgoing topic is empty")
		}

		if sensor.LowerBound >= sensor.HigherBound {
			return fmt.Errorf("lower bound is greater or equal to higher bound")
		}
	}

	return nil
}

func LoadDefaultAlertConfig() (*AlertConfig, error) {
	var cfg AlertConfig
	err := LoadDefaultConfig(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
