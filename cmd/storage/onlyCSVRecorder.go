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
				Token:        "gRGYqt6oxFui_8TBpMjkT7Vfq2XAz2wL48ldxzYHp-cADW14GmN0dYYIMxtqDBcZIgk5YdgvZn_rrEzHhjYw1w==",
				Bucket:       "metrics",
				Organization: "meteo-airport",
				Measurement:  "temperature",
			}); err != nil {
			panic(err)
		} else {
			manager.AddRecorder(sensor.Temperature, influxRecorder, 1)
		}

		if err := manager.Start(); err != nil {
			panic(err)
		}

		defer func(manager *storage.Manager) {
			err := manager.Close()
			if err != nil {
				panic(err)
			}
		}(manager)

		measurement := sensor.Measurement{
			SensorID:  1,
			AirportID: "NTE",
			Value:     21.0,
			Unit:      "°C",
			Timestamp: time.Now(),
		}

		if err := measurement.PublishOnMQTT(sensor.Temperature, 1, false, client); err != nil {
			panic(err)
		}

	}

}
