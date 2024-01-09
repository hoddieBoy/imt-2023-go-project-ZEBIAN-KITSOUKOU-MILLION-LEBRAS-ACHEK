package mqtt_helper

import (
	"fmt"
	"imt-atlantique.project.group.fr/meteo-airport/internal/logutil"
	"os"
)

type MQTTConfig struct {
	Server struct {
		Protocol string `yaml:"protocol"`
		Port     int    `yaml:"port"`
		Host     string `yaml:"host"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"server"`
}

type SensorAlertType struct {
	EndPoint    string `yaml:"end_point"`
	LowerBound  int    `yaml:"lower_bound"`
	HigherBound int    `yaml:"higher_bound"`
}
type MQTTRoot struct {
	Root struct {
		Sensor struct {
			Humidity    string `yaml:"humidity"`
			Temperature string `yaml:"temperature"`
			Pressure    string `yaml:"pressure"`
		} `yaml:"sensor"`

		Alert struct {
			Humidity    SensorAlertType `yaml:"humidity"`
			Temperature SensorAlertType `yaml:"temperature"`
			Pressure    SensorAlertType `yaml:"pressure"`
		} `yaml:"alert"`
	} `yaml:"root-topic"`
}

func RetrieveMQTTPropertiesFromYaml(filePath string) (*MQTTConfig, error) {
	f, err := os.Open(filePath)
	if err != nil {
		logutil.Error("Failed to open file:\n\t << %v >>", err)
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	var cfg MQTTConfig
	decoder := yaml.NewDecoder(f)
	if decoder.Decode(&cfg) != nil {
		logutil.Error("Failed to decode file: << %v >>", err)
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		logutil.Error("Failed to validate config: << %v >>", err)
		return nil, err
	}

	return &cfg, err
}

func RetrieveMQTTRootFromYaml() (*MQTTRoot, error) {
	f, err := os.Open("./config/message-topic.yaml")
	if err != nil {
		return nil, fmt.Errorf("\033[31mfailed to open file:\n\t<<%w>>\033[0m", err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	var cfg MQTTRoot
	decoder := yaml.NewDecoder(f)
	if decoder.Decode(&cfg) != nil {
		return nil, fmt.Errorf("\033[31mfailed to decode file:\n\t<<%w>>\033[0m", err)
	}

	return &cfg, err
}

func (c *MQTTConfig) GetServerAddress() string {
	return fmt.Sprintf("%s://%s:%d", c.Server.Protocol, c.Server.Host, c.Server.Port)
}

func (c *MQTTConfig) Validate() error {
	if c.Server.Host == "" {
		return fmt.Errorf("host is empty")
	}

	if c.Server.Port == 0 {
		return fmt.Errorf("port is empty")
	}

	if c.Server.Protocol == "" {
		return fmt.Errorf("protocol is empty")
	}

	return nil
}
