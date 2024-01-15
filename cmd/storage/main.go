package main

import (
	"math/rand"
	"time"

	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
	"imt-atlantique.project.group.fr/meteo-airport/internal/storage"
)

func createMQTTClient() *mqtt.Client {
	config, configErr := mqtt.RetrieveMQTTPropertiesFromYaml("./config/hiveClientConfig.yaml")
	if configErr != nil {
		panic(configErr)
	}

	client := mqtt.NewClient(config, "aClientId")

	if connexionErr := client.Connect(); connexionErr != nil {
		panic(connexionErr)
	}

	return client
}

func createManager(client *mqtt.Client) *storage.Manager {
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

	influxRecorder, err := storage.NewInfluxDBRecorder(
		storage.InfluxDBSettings{
			URL:          "http://localhost:8086",
			Token:        "hDwq6Hds2yXjMjDHCFjBNZZ_vOsEbF4DdKvUfnjb8rMNkTRjCrOwnoLfPf9Oy7eOqHsvawau36-DVqHUwvKNGw==",
			Bucket:       "metrics",
			Organization: "meteo-airport",
		})

	if err != nil {
		panic(err)
	}

	manager.AddRecorder(sensor.Temperature, influxRecorder, 1)

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
	client := createMQTTClient()
	defer client.Disconnect()

	manager := createManager(client)

	if err := manager.Start(); err != nil {
		panic(err)
	}

	defer func(manager *storage.Manager) {
		err := manager.Close()
		if err != nil {
			panic(err)
		}
	}(manager)

	publishMeasurements(client)
}
