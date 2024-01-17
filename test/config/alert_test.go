package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"imt-atlantique.project.group.fr/meteo-airport/internal/config"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
)

func createValidAlert() *config.Alert {
	return &config.Alert{
		Broker: mqtt.Config{
			Protocol: "mqtt",
			Port:     1883,
			Host:     "testClient",
			Username: "testUser",
			Password: "testPassword",
		},
		SensorsAlert: map[string]config.SensorAlert{
			"sensor1": {
				IncomingTopic: "incoming",
				OutgoingTopic: "outgoing",
				LowerBound:    10.0,
				HigherBound:   20.0,
				ClientID:      "client1",
			},
		},
	}
}

func TestAlertValidationWithValidData(t *testing.T) {
	alert := createValidAlert()

	err := alert.Validate()

	assert.NoError(t, err)
}

func TestAlertValidationWithEmptyIncomingTopic(t *testing.T) {
	alert := createValidAlert()
	sensor := alert.SensorsAlert["sensor1"]
	sensor.IncomingTopic = ""

	err := alert.Validate()

	assert.Error(t, err)
	assert.Equal(t, "incoming topic is empty", err.Error())
}

func TestAlertValidationWithEmptyOutgoingTopic(t *testing.T) {
	alert := createValidAlert()
	sensor := alert.SensorsAlert["sensor1"]
	sensor.OutgoingTopic = ""

	err := alert.Validate()

	assert.Error(t, err)
	assert.Equal(t, "outgoing topic is empty", err.Error())
}

func TestAlertValidationWithInvalidBounds(t *testing.T) {
	alert := createValidAlert()
	sensor := alert.SensorsAlert["sensor1"]
	sensor.LowerBound = 20.0
	sensor.HigherBound = 10.0

	err := alert.Validate()

	assert.Error(t, err)
	assert.Equal(t, "lower bound is greater or equal to higher bound", err.Error())
}

func TestAlertValidationWithEmptyClientID(t *testing.T) {
	alert := createValidAlert()
	sensor := alert.SensorsAlert["sensor1"]
	sensor.ClientID = ""

	err := alert.Validate()

	assert.Error(t, err)
	assert.Equal(t, "client id is empty", err.Error())
}
