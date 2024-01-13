package main

import (
	"imt-atlantique.project.group.fr/meteo-airport/internal/config_helper"
	"imt-atlantique.project.group.fr/meteo-airport/internal/logutil"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt_helper"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
	"imt-atlantique.project.group.fr/meteo-airport/internal/storage"
	"os"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		logutil.Warn("No config file specified, using default one")
	}
	config_helper.SetDefaultConfigFileName("storage_config.yaml")
	if config, err := config_helper.LoadDefaultStorageConfig(); err != nil {
		panic(err)
	} else {
		client := mqtt_helper.NewClient(&config.MQTT, "aClientId")
		if err := client.Connect(); err != nil {
			panic(err)
		}
		defer client.Disconnect()
		logutil.Info("Config: %v", config)
		manager := storage.NewManager(client)

		if err := manager.Start(); err != nil {
			panic(err)
		}

		measurement := sensor.Measurement{
			SensorID:  1,
			AirportID: "NTE",
			Type:      sensor.Temperature,
			Value:     40.0,
			Unit:      "°C",
			Timestamp: time.Now(),
		}

		if err := measurement.PublishOnMQTT(1, false, client); err != nil {
			panic(err)
		}

	}

}
