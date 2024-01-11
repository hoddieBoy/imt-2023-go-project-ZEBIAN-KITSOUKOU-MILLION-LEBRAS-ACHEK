package config_helper

import (
	"fmt"
	"imt-atlantique.project.group.fr/meteo-airport/internal/logutil"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt_helper"
)

type AlertConfig struct {
	Broker       *mqtt_helper.MQTTConfig `yaml:"broker"`
	SensorsAlert []SensorAlert           `yaml:"sensors_alert"`
}

type SensorAlert struct {
	SensorType    string `yaml:"sensor_type"`
	IncomingTopic string `yaml:"incoming_topic"`
	OutgoingTopic string `yaml:"outgoing_topic"`
	LowerBound    int    `yaml:"lower_bound"`
	HigherBound   int    `yaml:"higher_bound"`
}

func (c *AlertConfig) Validate() error {
	// Check if the MQTTConfig is valid
	if err := c.Broker.Validate(); err != nil {
		logutil.Error("Invalid MQTT configuration")
		return err
	}

	// Check if the incoming topic is valid
	for _, sensor := range c.SensorsAlert {
		if sensor.IncomingTopic == "" {
			logutil.Error("Invalid incoming topic")
			return fmt.Errorf("incoming topic is empty")
		}

		if sensor.OutgoingTopic == "" {
			logutil.Error("Invalid outgoing topic")
			return fmt.Errorf("outgoing topic is empty")
		}

		if sensor.LowerBound >= sensor.HigherBound {
			logutil.Error("Invalid bounds")
			return fmt.Errorf("lower bound is greater or equal to higher bound")
		}
	}

	return nil
}

func RetrieveAlertPropertiesFromYaml(filePath string) (*AlertConfig, error) {
	var cfg AlertConfig
	err := RetrievePropertiesFromYaml(filePath, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func LoadDefaultAlertConfig() (*AlertConfig, error) {
	var cfg AlertConfig
	err := LoadDefaultConfig(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
