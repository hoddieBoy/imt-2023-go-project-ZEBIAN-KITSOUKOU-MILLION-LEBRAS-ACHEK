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

	humiditySensorConfig, configErr := config.LoadDefaultSensorConfig()

	if configErr != nil {
		log.Error("Error loading humiditySensorConfig: %v", configErr)
		os.Exit(1)
	}

	humiditySensor, err := sensor.InitializeSensor(humiditySensorConfig)

	if err != nil {
		panic(err)
	}

	humidity := 50.0
	minimalValue := 40.0
	maximalValue := 60.0

	for {
		currentHumidity := readHumidity(humidity, minimalValue, maximalValue)
		humiditySensor.UpdateLastMeasurement(currentHumidity)
		err := humiditySensor.PublishData()

		if err != nil {
			log.Error(fmt.Sprintf("Failed to publish data to client: %v", err))
		}

		fmt.Printf("Humidity: %f\n", humidity)
		time.Sleep(4 * time.Second)
	}
}

func readHumidity(currentHumidity float64, min float64, max float64) float64 {
	// Simulating humidity between 40 and 60%
	simulatedHumidity := currentHumidity

	if simulatedHumidity < min {
		simulatedHumidity = min
	}

	if simulatedHumidity > max {
		simulatedHumidity = max
	}

	simulatedHumidity += rand.Float64() * 2

	return simulatedHumidity
}
