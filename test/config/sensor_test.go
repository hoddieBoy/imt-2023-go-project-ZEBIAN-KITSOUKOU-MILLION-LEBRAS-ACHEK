package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"imt-atlantique.project.group.fr/meteo-airport/internal/config"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
)

func createValidSensorConfig() *config.SensorConfig {
	return &config.SensorConfig{
		Setting: struct {
			AirportID  string `yaml:"airport_id"`
			SensorID   int64  `yaml:"sensor_id"`
			Topic      string `yaml:"topic"`
			Type       string `yaml:"type"`
			Unit       string `yaml:"unit"`
			TimeFormat string `yaml:"time_format"`
		}{
			AirportID:  "airport1",
			SensorID:   1,
			Topic:      "topic1",
			Type:       "type1",
			Unit:       "unit1",
			TimeFormat: "time_format1",
		},
		Broker: struct {
			Client   mqtt.Config `yaml:"client"`
			Qos      byte        `yaml:"qos"`
			ClientID string      `yaml:"client_id"`
		}{
			Client: mqtt.Config{
				Protocol: "mqtt",
				Port:     1883,
				Host:     "testClient",
				Username: "testUser",
				Password: "testPassword",
			},
			Qos:      1,
			ClientID: "client1",
		},
	}
}

func TestSensorConfigValidationWithValidData(t *testing.T) {
	sensorConfig := createValidSensorConfig()

	err := sensorConfig.Validate()

	assert.NoError(t, err)
}

func TestSensorConfigValidationWithEmptyAirportID(t *testing.T) {
	sensorConfig := createValidSensorConfig()
	sensorConfig.Setting.AirportID = ""

	err := sensorConfig.Validate()

	assert.Error(t, err)
	assert.Equal(t, "airport id is empty", err.Error())
}

func TestSensorConfigValidationWithInvalidQos(t *testing.T) {
	sensorConfig := createValidSensorConfig()
	sensorConfig.Broker.Qos = 3

	err := sensorConfig.Validate()

	assert.Error(t, err)
	assert.Equal(t, "qos must be between 0 and 2", err.Error())
}
