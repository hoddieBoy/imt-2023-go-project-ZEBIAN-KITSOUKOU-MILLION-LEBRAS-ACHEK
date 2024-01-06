package main

import (
	"github.com/eclipse/paho.mqtt.golang"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt-helper"
)

func main() {
	if config, err := mqtt_helper.RetrievePropertiesFromConfig("./config/brokerConfig.yaml"); err != nil {
		panic(err)
	} else {
		client := mqtt_helper.NewClient(config, "aClientId")
		if err := client.Connect(); err != nil {
			panic(err)
		}
		// Print the message when it is received
		if err := client.Subscribe("test", 1, func(client mqtt.Client, message mqtt.Message) {
			println(string(message.Payload()))
		}); err != nil {
			panic(err)
		}

		if err := client.Publish("test", 1, false, "Hello World"); err != nil {
			panic(err)
		}
		client.Disconnect()
	}
}
