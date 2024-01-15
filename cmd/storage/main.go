package main

import (
	"math/rand"
	"os"
	"time"

	"imt-atlantique.project.group.fr/meteo-airport/internal/config"
	"imt-atlantique.project.group.fr/meteo-airport/internal/log"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
	"imt-atlantique.project.group.fr/meteo-airport/internal/storage"
)

func createMQTTClient(config *config.Storage) *mqtt.Client {
	client := mqtt.NewClient(&config.MQTT)

	if connexionErr := client.Connect(); connexionErr != nil {
		log.Error("Error connecting to MQTT broker: %v", connexionErr)
		os.Exit(1)
	}

	return client
}

func createManager(storageConfig *config.Storage) *storage.Manager {
	client := mqtt.NewClient(&storageConfig.MQTT)

	if connexionErr := client.Connect(); connexionErr != nil {
		log.Error("Error connecting to MQTT broker: %v", connexionErr)
		os.Exit(1)
	}

	manager := storage.NewManager(client)

	for measurement, storageConfigs := range storageConfig.Settings {
		if storageConfigs.InfluxDB != (config.InfluxDBSettings{}) {
			log.Info("Registering InfluxDB recorder for measurement %s", measurement)

			influxDBRecorder, _ := storage.NewInfluxDBRecorder(storageConfigs.InfluxDB)
			manager.AddRecorder(sensor.MeasurementType(measurement), influxDBRecorder, 1)
		}

		if storageConfigs.CSV != (config.CSVSettings{}) {
			log.Info("Registering CSV recorder for measurement %s", measurement)

			csvRecorder, _ := storage.NewCSVRecorder(storageConfigs.CSV)
			manager.AddRecorder(sensor.MeasurementType(measurement), csvRecorder, 1)
		}
	}

	return manager
}

func publishMeasurements(client *mqtt.Client) {
	measurement := sensor.Measurement{
		SensorID:  1,
		AirportID: "NTE",
		Type:      sensor.Temperature,
		Value:     20.0,
		Unit:      "°C",
		Timestamp: time.Now(),
	}

	for {
		measurement.Timestamp = time.Now()
		measurement.Value = measurement.Value + rand.Float64() - 0.5

		if err := measurement.PublishOnMQTT(1, false, client); err != nil {
			panic(err)
		}

		time.Sleep(1 * time.Second)
	}
}

func main() {
	log.Info("Loading defaultStorageConfig...")

	defaultStorageConfig, configErr := config.LoadDefaultStorageConfig()

	if configErr != nil {
		log.Error("Error loading defaultStorageConfig: %v", configErr)
		os.Exit(1)
	}

	log.Info("Starting storage manager...")

	manager := createManager(defaultStorageConfig)

	if err := manager.Start(); err != nil {
		log.Error("Error starting storage manager: %v", err)
		os.Exit(1)
	}

	defer func(manager *storage.Manager) {
		err := manager.Close()
		if err != nil {
			log.Error("Error closing storage manager: %v", err)
		}
	}(manager)

	client := createMQTTClient(defaultStorageConfig)

	publishMeasurements(client)
}
