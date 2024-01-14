package mqtt

import (
	"fmt"
)

type MQTTConfig struct {
	Protocol string
	Port     int
	Host     string
	Username string
	Password string
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s://%s:%d", c.Protocol, c.Host, c.Port)
}

func (c *Config) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("host is empty")
	}

	if c.Port == 0 {
		return fmt.Errorf("port is empty")
	}

	if c.Protocol == "" {
		return fmt.Errorf("protocol is empty")
	}

	return nil
}
