package sensor

type MeasurementType string

const (
	Temperature MeasurementType = "temperature"
	Humidity    MeasurementType = "humidity"
	Pressure    MeasurementType = "pressure"
	WindSpeed   MeasurementType = "windSpeed"
)

func (t MeasurementType) GetTopic() string {
	return "airport/+/" + string(t)
}
