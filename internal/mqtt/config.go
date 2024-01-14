package mqtt

import (
	"fmt"
)

type Config struct {
	Protocol string `yaml:"protocol"`
	Port     int    `yaml:"port"`
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	ClientID string `yaml:"client_id"`
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

	if c.ClientID == "" {
		return fmt.Errorf("clientID is empty")
	}

	return nil
}
