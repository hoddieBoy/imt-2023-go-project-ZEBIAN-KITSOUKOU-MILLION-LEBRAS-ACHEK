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

func publishData(actualPressure float64, sensor sensor.Sensor) {
	sensor.GenerateData(2, "CGD", "pressure",
		actualPressure, "hPa", time.Now())

	err := sensor.PublishData()

	if err != nil {
		log.Error(fmt.Sprintf("Failed to publish data to client: %v", err))
	}
}

func main() {
	sensor := sensor.Sensor{}
	err := sensor.InitializeSensor()

	if err != nil {
		panic(err)
	}

	actualPressure := 1013.25
	minimalValue := 10.0
	maximalValue := 1083.8

	for {
		actualPressure = windDataGeneration(actualPressure, minimalValue, maximalValue)
		publishData(actualPressure, sensor)
		time.Sleep(5 * time.Second)
	}
}
