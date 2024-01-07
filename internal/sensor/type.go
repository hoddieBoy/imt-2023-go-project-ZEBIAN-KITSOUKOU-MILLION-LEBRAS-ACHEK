package sensor

import "time"

type Type string

const (
	Temperature Type = "temperature"
	Humidity    Type = "humidity"
	Pressure    Type = "pressure"
	WindSpeed   Type = "windSpeed"
)

func (t Type) GetTopic() string {
	return "airport/+/" + time.Now().Format("2006-01-02") + "/" + string(t)
}
