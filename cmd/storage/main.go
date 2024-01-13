package main

import (
	"math/rand"
	"time"

	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
	"imt-atlantique.project.group.fr/meteo-airport/internal/storage"
)

func main() {
	config, configErr := mqtt.RetrieveMQTTPropertiesFromYaml("./config/hiveClientConfig.yaml")
	if configErr != nil {
		panic(configErr)
	}

	client := mqtt.NewClient(config, "aClientId")

	if connexionErr := client.Connect(); connexionErr != nil {
		panic(connexionErr)
	}

	defer client.Disconnect()

	manager := storage.NewManager(client)
	csvSettings := storage.CSVSettings{
		PathDirectory: "./data",
		Separator:     ';',
		TimeFormat:    "2006-01-02 15:04:05",
	}

	csvRecorder, err := storage.NewCSVRecorder("test.csv", csvSettings)
	if err != nil {
		panic(err)
	}

	manager.AddRecorder(sensor.Temperature, csvRecorder, 1)

	if influxRecorder, err := storage.NewInfluxDBRecorder(
		storage.InfluxDBSettings{
			URL:          "http://localhost:8086",
			Token:        "hDwq6Hds2yXjMjDHCFjBNZZ_vOsEbF4DdKvUfnjb8rMNkTRjCrOwnoLfPf9Oy7eOqHsvawau36-DVqHUwvKNGw==",
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

	defer func(manager *storage.Manager) {
		err := manager.Close()
		if err != nil {
			panic(err)
		}
	}(manager)

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
