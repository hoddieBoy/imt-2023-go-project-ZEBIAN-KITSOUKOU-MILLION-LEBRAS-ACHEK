package main

import (
	"fmt"
	"math/rand"
	"time"

	"imt-atlantique.project.group.fr/meteo-airport/internal/log"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
)

func windDataGeneration(actualWind float64, min float64, max float64) float64 {
	actualWind += (rand.Float64() - rand.Float64()) * 10

	if actualWind < min {
		actualWind = min
	}

	if actualWind > max {
		actualWind = max
	}

	return actualWind
}

func publishData(sensor sensor.Sensor) {
	err := sensor.PublishData()

	if err != nil {
		log.Error(fmt.Sprintf("Failed to publish data to client: %v", err))
	}
}

func main() {
	actualWind := 40.0
	minimalValue := 10.0
	maximalValue := 120.0

	sensor := sensor.Sensor{}
	err := sensor.InitializeSensor(2, "CGD", "windSpeed",
		actualWind, "Km/h", time.Now())

	if err != nil {
		panic(err)
	}

	for {
		actualWind = windDataGeneration(actualWind, minimalValue, maximalValue)
		sensor.ChangeValueMeasurement(actualWind)
		publishData(sensor)
		time.Sleep(5 * time.Second)
	}
}
