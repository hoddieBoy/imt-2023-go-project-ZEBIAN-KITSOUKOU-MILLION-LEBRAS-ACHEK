package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"imt-atlantique.project.group.fr/meteo-airport/internal/config"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
)

func createValidStorage() *config.Storage {
	return &config.Storage{
		Settings: map[string]config.Setting{
			"setting1": {
				InfluxDB: config.InfluxDBSettings{
					URL:          "http://localhost:8086",
					Token:        "token",
					Bucket:       "bucket",
					Organization: "org",
				},
				CSV: config.CSVSettings{
					PathDirectory: "/path/to/directory",
					Separator:     ",",
					TimeFormat:    "2006-01-02T15:04:05Z07:00",
				},
				Qos:   1,
				Topic: "topic",
			},
		},
		Broker: struct {
			Config   mqtt.Config `yaml:"config"`
			ClientID string      `yaml:"client_id"`
		}{
			Config: mqtt.Config{
				Protocol: "mqtt",
				Port:     1883,
				Host:     "testClient",
				Username: "testUser",
				Password: "testPassword",
			},
			ClientID: "client1",
		},
	}
}

func TestStorageValidationWithValidData(t *testing.T) {
	storage := createValidStorage()

	err := storage.Validate()

	assert.NoError(t, err)
}

func TestStorageValidationWithEmptySettings(t *testing.T) {
	storage := createValidStorage()
	storage.Settings = map[string]config.Setting{}

	err := storage.Validate()

	assert.Error(t, err)
	assert.Equal(t, "settings is empty", err.Error())
}

func TestStorageValidationWithInvalidQos(t *testing.T) {
	storage := createValidStorage()
	setting := storage.Settings["setting1"]
	setting.Qos = 3

	err := storage.Validate()

	assert.Error(t, err)
	assert.Equal(t, "qos must be between 0 and 2", err.Error())
}

func TestStorageValidationWithEmptyClientID(t *testing.T) {
	storage := createValidStorage()
	storage.Broker.ClientID = ""

	err := storage.Validate()

	assert.Error(t, err)
	assert.Equal(t, "client id is empty", err.Error())
}
