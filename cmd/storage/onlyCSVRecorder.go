package main

import (
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt_helper"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
	"imt-atlantique.project.group.fr/meteo-airport/internal/storage"
	"time"
)

func main() {
	if config, err := mqtt_helper.RetrievePropertiesFromConfig("./config/brokerLocalhostConfig.yaml"); err != nil {
		panic(err)
	} else {
		client := mqtt_helper.NewClient(config, "aClientId")
		if err := client.Connect(); err != nil {
			panic(err)
		}
		defer client.Disconnect()

		measurement := &sensor.Measurement{
			SensorID:  1,
			AirportID: "NTE",
			Value:     20.0,
			Unit:      "°C",
			Timestamp: time.Now(),
		}

		manager := storage.NewManager(client)
		csvSettings := storage.CSVSettings{
			PathDirectory: "./data",
			Separator:     ';',
			TimeFormat:    "2006-01-02 15:04:05",
		}

		if csvRecorder, err := storage.NewCSVRecorder("test.csv", csvSettings); err != nil {
			panic(err)
		} else {
			manager.AddRecorder(sensor.Temperature, csvRecorder)
		}

		manager.Start()

		if err := sensor.PublishMeasurement(measurement, "temperature", 1, false, client); err != nil {
			panic(err)
		}
	}

}
