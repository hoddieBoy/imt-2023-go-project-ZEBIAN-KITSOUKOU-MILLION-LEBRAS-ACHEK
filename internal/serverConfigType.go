package internal

import "gopkg.in/yaml.v2"
import "os"

type Config struct {
	Server struct {
		Port     string `yaml:"port"`
		Host     string `yaml:"host"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"server"`
}

func retrievePropertiesFromConfig() *Config {
	f, err := os.Open("serverConfigType.yml")
	if err != nil {
		return nil
	}
	defer f.Close()
	var cfg Config
	decoder := yaml.NewDecoder(f)
	if decoder.Decode(&cfg) != nil {
		return nil
	}
	return &cfg
}
