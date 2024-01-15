package mqtt

import (
	"fmt"
<<<<<<< HEAD:internal/mqtt_helper/config.go
	"gopkg.in/yaml.v3"
	"imt-atlantique.project.group.fr/meteo-airport/internal/logutil"
=======
>>>>>>> main:internal/mqtt/config.go
	"os"

	"gopkg.in/yaml.v3"
	"imt-atlantique.project.group.fr/meteo-airport/internal/log"
)

type Config struct {
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
type Root struct {
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

func RetrieveMQTTPropertiesFromYaml(filePath string) (*Config, error) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Error("Failed to open file:\n\t << %v >>", err)
		return nil, err
	}

	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	var cfg Config

	decoder := yaml.NewDecoder(f)

	if decoder.Decode(&cfg) != nil {
		log.Error("Failed to decode file: << %v >>", err)
		return nil, err
	}

	if validationErr := cfg.Validate(); validationErr != nil {
		log.Error("Failed to validate config: << %v >>", validationErr)
		return nil, err
	}

	return &cfg, err
}

func RetrieveMQTTRootFromYaml() (*Root, error) {
	f, err := os.Open("./config/message-topic.yaml")
	if err != nil {
		return nil, fmt.Errorf("\033[31mfailed to open file:\n\t<<%w>>\033[0m", err)
	}

	defer func() {
		if closeErr := f.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}()

	var cfg Root

	decoder := yaml.NewDecoder(f)

	if decoder.Decode(&cfg) != nil {
		return nil, fmt.Errorf("\033[31mfailed to decode file:\n\t<<%w>>\033[0m", err)
	}

	return &cfg, err
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s://%s:%d", c.Server.Protocol, c.Server.Host, c.Server.Port)
}

func (c *Config) Validate() error {
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
