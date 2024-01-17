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
	actualWind += (rand.Float64() - rand.Float64()) * 10

	if actualWind < min {
		actualWind = min
	}

	if actualWind > max {
		actualWind = max
	}

	return actualWind
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		log.Warn("No config file specified, using default path: config/config.yaml")
	} else {
		config.SetDefaultConfigFileName(args[0])
	}

	windSensorConfig, configErr := config.LoadDefaultSensorConfig()

	if configErr != nil {
		log.Error("Error loading windSensorConfig: %v", configErr)
		os.Exit(1)
	}

	actualWind := 40.0
	minimalValue := 10.0
	maximalValue := 120.0

	windSensor, err := sensor.InitializeSensor(windSensorConfig)

	if err != nil {
		panic(err)
	}

	for {
		actualWind = windDataGeneration(actualWind, minimalValue, maximalValue)
		windSensor.UpdateLastMeasurement(actualWind)
		err := windSensor.PublishData()

		if err != nil {
			log.Error(fmt.Sprintf("Failed to publish data to client: %v", err))
		}

		time.Sleep(5 * time.Second)
	}
}
