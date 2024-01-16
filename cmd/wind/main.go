package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"imt-atlantique.project.group.fr/meteo-airport/internal/config"
	"imt-atlantique.project.group.fr/meteo-airport/internal/log"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
)

func windDataGeneration(actualWind float64, min float64, max float64) float64 {
	actualWind += (rand.Float64() - rand.Float64()) * 5

	if actualWind < min {
		actualWind = min
	}

	if actualWind > max {
		actualWind = max
	}

	return actualWind
}

func publishData(actualWind float64, sensor *sensor.Sensor) {
	sensor.GenerateData(actualWind)

	err := sensor.PublishData()

	if err != nil {
		log.Error(fmt.Sprintf("Failed to publish data to client: %v", err))
	}
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		log.Warn("No config file specified, using default path: config/config.yaml")
	} else {
		config.SetDefaultConfigFileName(args[0])
	}

	windSensor, err := sensor.InitializeSensor(sensor.WindSpeed)

	if err != nil {
		panic(err)
	}

	actualWind := 40.0
	minimalValue := 10.0
	maximalValue := 120.0

	for {
		actualWind = windDataGeneration(actualWind, minimalValue, maximalValue)
		publishData(actualWind, windSensor)
		time.Sleep(5 * time.Second)
	}
}
