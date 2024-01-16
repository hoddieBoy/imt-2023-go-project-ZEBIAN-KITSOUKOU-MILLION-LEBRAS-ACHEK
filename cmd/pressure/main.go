package main

import (
	"fmt"
	"math/rand"
	"time"

	"imt-atlantique.project.group.fr/meteo-airport/internal/log"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
)

func windDataGeneration(actualPressure float64, min float64, max float64) float64 {
	actualPressure += rand.Float64() - rand.Float64()

	if actualPressure < min {
		actualPressure = min
	}

	if actualPressure > max {
		actualPressure = max
	}

	return actualPressure
}

func publishData(sensor sensor.Sensor) {

	err := sensor.PublishData()

	if err != nil {
		log.Error(fmt.Sprintf("Failed to publish data to client: %v", err))
	}
}

func main() {
	actualPressure := 1013.25
	minimalValue := 875.0
	maximalValue := 1083.8

	sensor := sensor.Sensor{}
	err := sensor.InitializeSensor(2, "CGD", "windSpeed",
		actualPressure, "Km/h", time.Now())

	if err != nil {
		panic(err)
	}

	for {
		actualPressure = windDataGeneration(actualPressure, minimalValue, maximalValue)
		sensor.ChangeValueMeasurement(actualPressure)
		publishData(sensor)
		time.Sleep(5 * time.Second)
	}
}
