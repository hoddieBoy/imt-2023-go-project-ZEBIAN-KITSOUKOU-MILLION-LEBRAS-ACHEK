package sensor

type MeasurementType string

const (
	Temperature MeasurementType = "temperature"
	Humidity    MeasurementType = "humidity"
	Pressure    MeasurementType = "pressure"
	WindSpeed   MeasurementType = "windSpeed"
)
