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

		manager := storage.NewManager(client)
		csvSettings := storage.CSVSettings{
			PathDirectory: "./data",
			Separator:     ';',
			TimeFormat:    "2006-01-02 15:04:05",
		}

		if csvRecorder, err := storage.NewCSVRecorder("test.csv", csvSettings); err != nil {
			panic(err)
		} else {
			manager.AddRecorder(sensor.Temperature, csvRecorder, 1)
		}

		if influxRecorder, err := storage.NewInfluxDBRecorder(
			storage.InfluxDBSettings{
				URL:          "http://localhost:8086",
				Token:        "oovEfiXEVMcbEFEmy9jNNKUnWvlBexVmvSH6Ohr8f5b201RHNx5nb4D_5Y4o6mhwUWx--sV018H0ZLf6zDaD8Q==",
				Bucket:       "metrics",
				Organization: "meteo-airport",
			}); err != nil {
			panic(err)
		} else {
			manager.AddRecorder(sensor.Temperature, influxRecorder, 1)
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
