package sensor

import (
	"encoding/json"
	"fmt"
	"time"

	"imt-atlantique.project.group.fr/meteo-airport/internal/log"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
)

// Measurement represents the data from a sensor
type Measurement struct {
	SensorID  int64           `json:"sensor_id"`
	AirportID string          `json:"airport_id"`
	Type      MeasurementType `json:"type"`
	Value     float64         `json:"value"`
	Unit      string          `json:"unit"`
	Timestamp time.Time       `json:"timestamp"`
}

func (m *Measurement) String() string {
	return fmt.Sprintf(
		"SensorID: %d, AirportID: %s, Value: %f, Unit: %s, Timestamp: %s",
		m.SensorID, m.AirportID, m.Value, m.Unit, m.Timestamp.Format(time.RFC3339),
	)
}

func (m *Measurement) ToJSON() ([]byte, error) {
	payload, err := json.Marshal(m)
	if err != nil {
		log.Error("Failed to marshal measurement to JSON: %v", err)
		return nil, err
	}

	return payload, nil
}

func FromJSON(payload []byte) (*Measurement, error) {
	var measurement Measurement
	if err := json.Unmarshal(payload, &measurement); err != nil {
		log.Error("Failed to unmarshal measurement from JSON: %v", err)
		return nil, err
	}

	return &measurement, nil
}

func (m *Measurement) ToCSV(separator string, timeFormat string) string {
	return fmt.Sprintf(
		"%d%s%s%s%s%s%f%s%s%s%s",
		m.SensorID, separator, m.AirportID, separator, m.Type, separator,
		m.Value, separator, m.Unit, separator, m.Timestamp.Format(timeFormat),
	)
}

func MeasurementFieldNames(separator string) string {
	return fmt.Sprintf(
		"sensor_id%sairport_id%stype%svalue%sunit%stimestamp",
		separator, separator, separator, separator, separator,
	)
}

// PublishOnMQTT publishes a measurement to the MQTT broker
func (m *Measurement) PublishOnMQTT(qos byte, retained bool, client *mqtt.Client) error {
	// Topic: airport/<airport_id>/<year-month-day>/<type_of_measurement>
	topic := fmt.Sprintf("airport/%s/%s/%s", m.AirportID, m.Timestamp.Format("2006-01-02"), m.Type)
	payload, err := m.ToJSON()

	if err != nil {
		return err
	}

	if err := client.Publish(topic, qos, retained, payload); err != nil {
		return err
	}

	return nil
}
