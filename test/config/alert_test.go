package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"imt-atlantique.project.group.fr/meteo-airport/internal/config"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
)

func TestAlertValidationWithValidData(t *testing.T) {
	alert := &config.Alert{
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

	err := alert.Validate()

	assert.NoError(t, err)
}

func TestAlertValidationWithEmptyIncomingTopic(t *testing.T) {
	alert := &config.Alert{
		Broker: mqtt.Config{
			Protocol: "mqtt",
			Port:     1883,
			Host:     "testClient",
			Username: "testUser",
			Password: "testPassword",
		},
		SensorsAlert: map[string]config.SensorAlert{
			"sensor1": {
				IncomingTopic: "",
				OutgoingTopic: "outgoing",
				LowerBound:    10.0,
				HigherBound:   20.0,
				ClientID:      "client1",
			},
		},
	}

	err := alert.Validate()

	assert.Error(t, err)
	assert.Equal(t, "incoming topic is empty", err.Error())
}

func TestAlertValidationWithEmptyOutgoingTopic(t *testing.T) {
	alert := &config.Alert{
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
				OutgoingTopic: "",
				LowerBound:    10.0,
				HigherBound:   20.0,
				ClientID:      "client1",
			},
		},
	}

	err := alert.Validate()

	assert.Error(t, err)
	assert.Equal(t, "outgoing topic is empty", err.Error())
}

func TestAlertValidationWithInvalidBounds(t *testing.T) {
	alert := &config.Alert{
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
				LowerBound:    20.0,
				HigherBound:   10.0,
				ClientID:      "client1",
			},
		},
	}

	err := alert.Validate()

	assert.Error(t, err)
	assert.Equal(t, "lower bound is greater or equal to higher bound", err.Error())
}

func TestAlertValidationWithEmptyClientID(t *testing.T) {
	alert := &config.Alert{
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
				ClientID:      "",
			},
		},
	}

	err := alert.Validate()

	assert.Error(t, err)
	assert.Equal(t, "client id is empty", err.Error())
}
