package sensor

import (
	"fmt"
	"imt-atlantique.project.group.fr/meteo-airport/internal/logutil"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt_helper"
	"time"
)

type Sensor struct {
	client *mqtt_helper.MQTTClient
	data   Measurement
}

func (s *Sensor) InitializeSensor() error {
	if config, err := mqtt_helper.RetrieveMQTTPropertiesFromYaml("./config/hiveClientConfig.yaml"); err != nil {
		panic(err)
	} else {
		client := mqtt_helper.NewClient(config, "clientId")

		err := client.Connect()
		if err != nil {
			logutil.Error(fmt.Sprintf("Cannot connect to client: %v", err))
			return err
		}

		s.client = client
		return nil
	}
}

func (s *Sensor) GenerateData(sensorId int64, airportId string, sensorType MeasurementType, value float64, unit string, timestamp time.Time) {

	s.data = Measurement{
		SensorID:  sensorId,
		AirportID: airportId,
		Type:      sensorType,
		Value:     value,
		Unit:      unit,
		Timestamp: timestamp,
	}

}

func (s *Sensor) PublishData() error {
	err := s.data.PublishOnMQTT(2, false, s.client)
	if err != nil {
		logutil.Error(fmt.Sprintf("Failed to publish data to client: %v", err))
		return err
	}

	return nil
}
