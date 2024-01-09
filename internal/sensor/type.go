package sensor

import "time"

type MeasurementType string

const (
	Temperature MeasurementType = "temperature"
	Humidity    MeasurementType = "humidity"
	Pressure    MeasurementType = "pressure"
	WindSpeed   MeasurementType = "windSpeed"
)

func (t MeasurementType) GetTopic() string {
	return "airport/+/" + time.Now().Format("2006-01-02") + "/" + string(t)
}
