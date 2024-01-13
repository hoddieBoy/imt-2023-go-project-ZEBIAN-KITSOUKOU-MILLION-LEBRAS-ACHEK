package sensor

import (
	"encoding/json"
	"fmt"
	"imt-atlantique.project.group.fr/meteo-airport/internal/logutil"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt_helper"
	"strconv"
	"strings"
	"time"
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
		logutil.Error("Failed to marshal measurement to JSON: %v", err)
		return nil, err
	}

	return payload, nil
}

func FromJSON(payload []byte) (*Measurement, error) {
	var measurement Measurement
	if err := json.Unmarshal(payload, &measurement); err != nil {
		logutil.Error("Failed to unmarshal measurement from JSON: %v", err)
		return nil, err
	}
	return &measurement, nil
}

func (m *Measurement) MeasurementValuesToCSV(separator rune, timeFormat string) string {
	data := []string{strconv.FormatInt(m.SensorID, 10),
		m.AirportID,
		string(m.Type),
		strconv.FormatFloat(m.Value, 'f', -1, 64),
		m.Unit,
		m.Timestamp.Format(timeFormat),
	}
	return strings.Join(data, strconv.QuoteRune(separator))
}

func MeasurementFieldNames(separator rune) string {
	return fmt.Sprintf(
		"sensor_id%cairport_id%ctype%cvalue%cunit%ctimestamp",
		separator, separator, separator, separator, separator,
	)
}

// PublishOnMQTT publishes a measurement to the MQTT broker
func (m *Measurement) PublishOnMQTT(qos byte, retained bool, client *mqtt_helper.MQTTClient) error {
	// Topic: airport/<airport_id>/<year-month-day>/<type_of_measurement>
	topic := fmt.Sprintf("airport/%s/%s/%s", m.AirportID, m.Timestamp.Format("2006-01-02"), m.Type)
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
