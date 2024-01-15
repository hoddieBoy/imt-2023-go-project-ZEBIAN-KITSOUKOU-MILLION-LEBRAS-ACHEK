package sensor

import (
	"fmt"
	"math/rand"
	"time"

	"imt-atlantique.project.group.fr/meteo-airport/internal/log"
)

func windDataGeneration(actualWind float64, min float64, max float64) float64 {

	actualWind = actualWind + (rand.Float64()-rand.Float64())*5

	if actualWind < min {
		actualWind = min
	}
	if actualWind > max {
		actualWind = max
	}

	return actualWind
}

func publishData(actualWind float64, sensor Sensor) {
	sensor.GenerateData(2, "CGD", MeasurementType("windSpeed"),
		actualWind, "Km/h", time.Now())
	err := sensor.PublishData()

	if err != nil {
		log.Error(fmt.Sprintf("Failed to publish data to client: %v", err))
	}
}

func main() {
	sensor := Sensor{}
	err := sensor.InitializeSensor()

	if err != nil {
		panic(err)
	}

	actualWind := 40.0
	minimalValue := 10.0
	maximalValue := 120.0

	for {
		actualWind := windDataGeneration(actualWind, minimalValue, maximalValue)
		publishData(actualWind, sensor)
		time.Sleep(5 * time.Second)
	}
}
