package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"imt-atlantique.project.group.fr/meteo-airport/internal/config_helper"
	"imt-atlantique.project.group.fr/meteo-airport/internal/logutil"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt_helper"
	"os"
)

func main() {
	if config, err := config_helper.LoadDefaultAlertConfig(); err != nil {
		panic(err)
	} else {
		client := mqtt_helper.NewClient(config.Broker, "anotherClientId")
		if err := client.Connect(); err != nil {
			panic(err)
		}
		defer client.Disconnect()

		handleAlertListening(client, config.SensorsAlert)

		select {}
	}
}

func handleAlertListening(client *mqtt_helper.MQTTClient, alerts []config_helper.SensorAlert) {

	for _, alert := range alerts {
		err := client.Subscribe(alert.IncomingTopic,
			1,
			checkValidRangeOnReception(client,
				alert,
				"Alert, "+alert.SensorType+" sensor out of range"))
		if err != nil {
			logutil.Error("failed to subscribe to topic %s:\n\t<<%v>>", alert.IncomingTopic, err)
			os.Exit(1)
		}
	}
}

func checkValidRangeOnReception(helperClient *mqtt_helper.MQTTClient, sensorAlert config_helper.SensorAlert, alertMessage string) func(client mqtt.Client, message mqtt.Message) {
	return func(mqttClient mqtt.Client, message mqtt.Message) {
		sensorValue := getJsonValueAsIntFromMessage(message)
		if !(sensorAlert.LowerBound <= sensorValue && sensorValue <= sensorAlert.HigherBound) {
			err := helperClient.Publish(sensorAlert.OutgoingTopic, 1, false, alertMessage)
			if err != nil {
				panic(err)
			}
		}
	}
}

func getJsonValueAsIntFromMessage(message mqtt.Message) int {
	//TODO
	return 50
}
