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

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		log.Warn("No config file specified, using default path: config/config.yaml")
	} else {
		config.SetDefaultConfigFileName(args[0])
	}

	pressureSensorConfig, configErr := config.LoadDefaultSensorConfig()

	if configErr != nil {
		log.Error("Error loading pressureSensorConfig: %v", configErr)
		os.Exit(1)
	}

	actualPressure := 1013.25
	minimalValue := 875.0
	maximalValue := 1083.8

	pressureSensor, err := sensor.InitializeSensor(pressureSensorConfig)

	if err != nil {
		log.Error("Error initializing sensor: %v", err)
		os.Exit(1)
	}

	for {
		actualPressure = windDataGeneration(actualPressure, minimalValue, maximalValue)
		pressureSensor.UpdateLastMeasurement(actualPressure)
		err := pressureSensor.PublishData()

		if err != nil {
			log.Error(fmt.Sprintf("Failed to publish data to client: %v", err))
		}

		time.Sleep(5 * time.Second)
	}
}
