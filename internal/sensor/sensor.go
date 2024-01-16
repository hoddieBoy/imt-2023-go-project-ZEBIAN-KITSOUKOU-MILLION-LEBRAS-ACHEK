package sensor

import (
	"fmt"
	"time"

	"imt-atlantique.project.group.fr/meteo-airport/internal/config"
	"imt-atlantique.project.group.fr/meteo-airport/internal/log"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
)

type Sensor struct {
	client       *mqtt.Client
	last         Measurement
	publishTopic string
	qos          byte
}

func InitializeSensor(config *config.SensorConfig) (*Sensor, error) {
	client := mqtt.NewClient(&config.Broker.Client)

	err := client.Connect()

	if err != nil {
		log.Error(fmt.Sprintf("Cannot connect to client: %v", err))
		return nil, err
	}

	sensor := &Sensor{
		client: client,
		last: Measurement{
			SensorID:  config.Setting.SensorID,
			AirportID: config.Setting.AirportID,
			Type:      MeasurementType(config.Setting.Type),
			Unit:      config.Setting.Unit,
		},
		publishTopic: config.Setting.Topic,
		qos:          config.Broker.Qos,
	}

	return sensor, nil
}

func (s *Sensor) PublishData() error {
	err := s.last.PublishOnMQTT(s.client, s.qos, false, s.publishTopic)

	if err != nil {
		log.Error(fmt.Sprintf("Failed to publish last to client: %v", err))
		return err
	}

	return nil
}

func (s *Sensor) UpdateLastMeasurement(value float64) {
	s.last.Value = value
	s.last.Timestamp = time.Now()
}
