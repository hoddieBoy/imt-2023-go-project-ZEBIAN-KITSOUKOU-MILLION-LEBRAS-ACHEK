package main

import (
	"fmt"

	pahoMqtt "github.com/eclipse/paho.mqtt.golang"
	"imt-atlantique.project.group.fr/meteo-airport/internal/config_helper"
	"imt-atlantique.project.group.fr/meteo-airport/internal/logutil"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
	"os"
)

func main() {
	if config, err := config_helper.LoadDefaultAlertConfig(); err != nil {
		panic(err)
	}

	client := mqtt.NewClient(config.Broker, "anotherClientId")
	if err := client.Connect(); err != nil {
		panic(err)
	}

	defer client.Disconnect()

		handleAlertListening(client, config.SensorsAlert)

		select {}
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

func checkValidRangeOnReception(
	helperClient *mqtt.Client,
	sensorAlert mqtt.SensorAlertType,
	alertMessage string,
) pahoMqtt.MessageHandler {
	return func(mqttClient pahoMqtt.Client, message pahoMqtt.Message) {
		sensorValue := getJSONValueAsIntFromMessage(message)
		if !(sensorAlert.LowerBound <= sensorValue && sensorValue <= sensorAlert.HigherBound) {
			err := helperClient.Publish(sensorAlert.OutgoingTopic, 1, false, alertMessage)
			if err != nil {
				panic(err)
			}
		}
	}
}

func getJSONValueAsIntFromMessage(message pahoMqtt.Message) int {
	// TODO
	fmt.Println(string(message.Payload()))
	return 50
}
