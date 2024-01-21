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

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		log.Warn("No config file specified, using default path: config/config.yaml")
	} else {
		config.SetDefaultConfigFileName(args[0])
	}

	tempSensorConfig, configErr := config.LoadDefaultSensorConfig()

	if configErr != nil {
		log.Error("Error loading tempSensorConfig: %v", configErr)
		os.Exit(1)
	}

	temperatureSensor, err := sensor.InitializeSensor(tempSensorConfig)

	if err != nil {
		panic(err)
	}

	temperature := 22.0
	minimalValue := 20.0
	maximalValue := 25.0

	for {
		currentTemperature := readTemp(temperature, minimalValue, maximalValue)
		temperatureSensor.UpdateLastMeasurement(currentTemperature)
		err := temperatureSensor.PublishData()

		if err != nil {
			log.Error(fmt.Sprintf("Failed to publish data to client: %v", err))
		}

		time.Sleep(3 * time.Second)
	}
}

func readTemp(currentTemperature float64, min float64, max float64) float64 {
	simulatedTemperature := currentTemperature

	if simulatedTemperature < min {
		simulatedTemperature = min
	}

	if simulatedTemperature > max {
		simulatedTemperature = max
	}

	simulatedTemperature += rand.Float64() * 2

	return simulatedTemperature
}
