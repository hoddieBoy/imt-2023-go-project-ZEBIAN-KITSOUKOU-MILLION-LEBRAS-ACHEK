package main

import (
	"math/rand"
	"os"

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

	log.Info("Loading Configurations of the sensor...")
	testSensor, err := sensor.InitializeSensor(sensor.Temperature)

	if err != nil {
		log.Error("Error loading defaultSensorConfig: %v", err)
		os.Exit(1)
	}

	log.Info("Starting sensor...")

	randomTemperature := 20.0
	for {

		if rand.Float64() > 0.5 {
			randomTemperature = randomTemperature + rand.Float64()
		} else {
			randomTemperature = randomTemperature - rand.Float64()
		}

		testSensor.GenerateData(randomTemperature)
		if err := testSensor.PublishData(); err != nil {
			log.Error("Error publishing data: %v", err)
			os.Exit(1)
		}
	}

}
