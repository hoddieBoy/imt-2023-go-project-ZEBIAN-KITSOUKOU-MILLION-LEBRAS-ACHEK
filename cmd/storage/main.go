package main

import (
	"imt-atlantique.project.group.fr/meteo-airport/internal/config_helper"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt_helper"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
	"imt-atlantique.project.group.fr/meteo-airport/internal/storage"
	"time"
)

func main() {
	config_helper.SetDefaultConfigFileName("storage_config.yaml")
	if config, err := config_helper.LoadDefaultStorageConfig(); err != nil {
		panic(err)
	} else {
		client := mqtt_helper.NewClient(config.MQTT, "aClientId")
		if err := client.Connect(); err != nil {
			panic(err)
		}
		defer client.Disconnect()

		manager := storage.NewManager(client)

		for measurement, storagesSettings := range config.Storages {
			for _, storageSettings := range storagesSettings {
				if storageSettings.InfluxDB != (config_helper.InfluxDBSettings{}) {
					if recorder, err := storage.NewInfluxDBRecorder(storageSettings.InfluxDB); err != nil {
						panic(err)
					} else {
						manager.AddRecorder(sensor.MeasurementType(measurement), recorder, 1)
					}
				}

				if storageSettings.CSV != (config_helper.CSVSettings{}) {
					if recorder, err := storage.NewCSVRecorder(storageSettings.CSV); err != nil {
						panic(err)
					} else {
						manager.AddRecorder(sensor.MeasurementType(measurement), recorder, 1)
					}
				}
			}
		}

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
