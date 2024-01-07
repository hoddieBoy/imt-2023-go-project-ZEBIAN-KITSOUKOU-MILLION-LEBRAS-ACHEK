package sensor

import (
	"encoding/json"
	"fmt"
	"imt-atlantique.project.group.fr/meteo-airport/internal/logutil"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt_helper"
	"strconv"
	"time"
)

// Measurement represents the data from a sensor
type Measurement struct {
	SensorID  int64     `json:"sensor_id"`
	AirportID string    `json:"airport_id"`
	Value     float64   `json:"value"`
	Unit      string    `json:"unit"`
	Timestamp time.Time `json:"timestamp"`
}

func (m *Measurement) String() string {
	return fmt.Sprintf(
		"SensorID: %d, AirportID: %s, Value: %f, Unit: %s, Timestamp: %s",
		m.SensorID, m.AirportID, m.Value, m.Unit, m.Timestamp,
	)
}

func (m *Measurement) ToJSON() ([]byte, error) {
	if payload, err := json.Marshal(m); err != nil {
		logutil.Error("Failed to marshal measurement to JSON: %v", err)
		return nil, err
	} else {
		return payload, nil
	}
}

func FromJSON(payload []byte) (*Measurement, error) {
	var measurement Measurement
	if err := json.Unmarshal(payload, &measurement); err != nil {
		logutil.Error("Failed to unmarshal measurement from JSON: %v", err)
		return nil, err
	}
	return &measurement, nil
}

func (m *Measurement) ToCSV(separator rune, timeFormat string) string {
	return fmt.Sprintf(
		"%d%c%s%c%s%c%s%c%s",
		m.SensorID, separator, m.AirportID, separator, strconv.FormatFloat(m.Value, 'f', -1, 64), separator, m.Unit, separator, m.Timestamp.Format(timeFormat),
	)
}

func MeasurementFieldNames(separator rune) string {
	return fmt.Sprintf(
		"sensor_id%cairport_id%cvalue%cunit%ctimestamp",
		separator, separator, separator, separator,
	)
}

// PublishOnMQTT publishes a measurement to the MQTT broker
func (m *Measurement) PublishOnMQTT(typeOfMeasurement Type, qos byte, retained bool, client *mqtt_helper.MQTTClient) error {
	// Topic: airport/<airport_id>/<year-month-day>/<type_of_measurement>
	topic := fmt.Sprintf("airport/%s/%s/%s", m.AirportID, m.Timestamp.Format("2006-01-02"), typeOfMeasurement)
	payload, err := m.ToJSON()
	if err != nil {
		logutil.Error(fmt.Sprintf("Failed to marshal measurement to JSON: %v", err))
		return err
	}

	if err := client.Publish(topic, qos, retained, payload); err != nil {
		logutil.Error(fmt.Sprintf("Failed to publish measurement to topic %s: %v", topic, err))
		return err
	}

	return nil
}
