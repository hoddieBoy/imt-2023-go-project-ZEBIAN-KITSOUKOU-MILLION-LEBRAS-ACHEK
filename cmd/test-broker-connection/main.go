package main

import (
	"github.com/eclipse/paho.mqtt.golang"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt_helper"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
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
		// Print the message when it is received
		if err := client.Subscribe("airport/#", 1, func(client mqtt.Client, message mqtt.Message) {
			println(string(message.Payload()))
		}); err != nil {
			panic(err)
		}

		if err := sensor.PublishMeasurement(measurement, "temperature", 1, false, client); err != nil {
			panic(err)
		}
	}
}
