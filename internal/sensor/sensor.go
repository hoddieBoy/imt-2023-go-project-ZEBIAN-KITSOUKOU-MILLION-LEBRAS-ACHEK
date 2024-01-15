package sensor

import (
	"fmt"
	"time"

	"imt-atlantique.project.group.fr/meteo-airport/internal/config"
	"imt-atlantique.project.group.fr/meteo-airport/internal/log"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
)

type Sensor struct {
	client *mqtt.Client
	config *config.SensorConfig
	last   *Measurement
}

func InitializeSensor(sensorType MeasurementType) (*Sensor, error) {
	cfg, err := config.LoadDefaultSensorConfig()

	if err != nil {
		return nil, err
	}

	client := mqtt.NewClient(&cfg.Broker.Client)

	if err := client.Connect(); err != nil {
		return nil, err
	}

	return &Sensor{
		client: client,
		config: cfg,
		last: &Measurement{
			SensorID:  cfg.Setting.SensorID,
			AirportID: cfg.Setting.AirportID,
			Type:      sensorType,
			Unit:      cfg.Setting.Unit,
		},
	}, nil
}

func (s *Sensor) GenerateData(value float64) {
	s.last.Value = value
	s.last.Timestamp = time.Now()
}

func (s *Sensor) PublishData() error {
	err := s.last.PublishOnMQTT(s.config.Broker.Qos, false, s.client)

	if err != nil {
		log.Error(fmt.Sprintf("Failed to publish data to client: %v", err))
		return err
	}

	return nil
}
