package config

import (
	"fmt"

	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
)

type Alert struct {
	Broker       mqtt.Config            `yaml:"broker"`
	SensorsAlert map[string]SensorAlert `yaml:"sensors_alert"`
}

type SensorAlert struct {
	IncomingTopic string `yaml:"incoming_topic"`
	OutgoingTopic string `yaml:"outgoing_topic"`
	LowerBound    int    `yaml:"lower_bound"`
	HigherBound   int    `yaml:"higher_bound"`
	IncomingQos   byte   `yaml:"incoming_qos"`
	OutgoingQos   byte   `yaml:"outgoing_qos"`
}

func (c *Alert) Validate() error {
	// Check if the MQTTConfig is valid
	if err := c.Broker.Validate(); err != nil {
		return err
	}

	for _, sensor := range c.SensorsAlert {
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

func LoadDefaultAlertConfig() (*Alert, error) {
	var cfg Alert
	err := LoadDefaultConfig(&cfg)

	return &cfg, err
}