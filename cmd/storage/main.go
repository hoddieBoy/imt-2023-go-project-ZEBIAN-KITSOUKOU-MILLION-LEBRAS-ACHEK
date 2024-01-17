package main

import (
	"os"

	"imt-atlantique.project.group.fr/meteo-airport/internal/config"
	"imt-atlantique.project.group.fr/meteo-airport/internal/log"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
	"imt-atlantique.project.group.fr/meteo-airport/internal/storage"
)

func createAndConnectClient(storageConfig *config.Storage) *mqtt.Client {
	client := mqtt.NewClient(&storageConfig.Broker.Config, storageConfig.Broker.ClientID)

	if connexionErr := client.Connect(); connexionErr != nil {
		log.Error("Error connecting to MQTT broker: %v", connexionErr)
		os.Exit(1)
	}

	return client
}

func registerInfluxDBRecorder(manager *storage.Manager, measurement string, setting config.Setting) {
	log.Info("Registering InfluxDB recorder for measurement %s", measurement)

	influxDBRecorder, err := storage.NewInfluxDBRecorder(setting.InfluxDB)
	if err != nil {
		log.Error("Error creating InfluxDB recorder: %v", err)
		os.Exit(1)
	}

	err = manager.AddRecorder(sensor.MeasurementType(measurement), setting.Topic, setting.Qos, influxDBRecorder)
	if err != nil {
		log.Error("Error adding recorder: %v", err)
		os.Exit(1)
	}
}

func registerCSVRecorder(manager *storage.Manager, measurement string, setting config.Setting) {
	log.Info("Registering CSV recorder for measurement %s", measurement)

	csvRecorder, _ := storage.NewCSVRecorder(setting.CSV)
	err := manager.AddRecorder(sensor.MeasurementType(measurement), setting.Topic, setting.Qos, csvRecorder)

	if err != nil {
		log.Error("Error adding recorder: %v", err)
		os.Exit(1)
	}

	if _, err := os.Stat(setting.CSV.PathDirectory); os.IsNotExist(err) {
		log.Info("Creating the directory for saving the CSV files...")

		err := os.Mkdir(setting.CSV.PathDirectory, 0755)

		if err != nil {
			log.Error("Error creating the directory for saving the CSV files: %v", err)
			os.Exit(1)
		}
	}
}

func createManager(storageConfig *config.Storage) *storage.Manager {
	client := createAndConnectClient(storageConfig)
	manager := storage.NewManager(client)

	for measurement, setting := range storageConfig.Settings {
		if setting.InfluxDB != (config.InfluxDBSettings{}) {
			registerInfluxDBRecorder(manager, measurement, setting)
		}

		if setting.CSV != (config.CSVSettings{}) {
			registerCSVRecorder(manager, measurement, setting)
		}
	}

	return manager
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		log.Warn("No config file specified, using default path: config/config.yaml")
	} else {
		config.SetDefaultConfigFileName(args[0])
	}

	log.Info("Loading Configurations of the storage manager...")

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

	select {}
}
