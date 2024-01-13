package main

import (
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

	manager.AddRecorder(sensor.Temperature, csvRecorder)

	manager.Start()
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

	if err := measurement.PublishOnMQTT(1, false, client); err != nil {
		panic(err)
	}
}
