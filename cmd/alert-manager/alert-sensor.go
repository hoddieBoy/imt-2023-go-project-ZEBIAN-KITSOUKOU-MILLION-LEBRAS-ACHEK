package main

import (
	"fmt"

	"os"

	pahoMqtt "github.com/eclipse/paho.mqtt.golang"
	"imt-atlantique.project.group.fr/meteo-airport/internal/config"
	"imt-atlantique.project.group.fr/meteo-airport/internal/log"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
)

func main() {
	alertConfig, err := config.LoadDefaultAlertConfig()

	if err != nil {
		panic(err)
	}

	client := mqtt.NewClient(&alertConfig.Broker)
	if err := client.Connect(); err != nil {
		panic(err)
	}

	defer client.Disconnect()

	handleAlertListening(client, alertConfig.SensorsAlert)

	select {}
}

func handleAlertListening(client *mqtt.Client, alerts map[string]config.SensorAlert) {
	for sensorType, alert := range alerts {
		err := client.Subscribe(alert.IncomingTopic,
			alert.IncomingQos,
			checkValidRangeOnReception(client,
				alert,
				"Alert, "+sensorType+" sensor out of range"))
		if err != nil {
			log.Error("failed to subscribe to topic %s:\n\t<<%v>>", alert.IncomingTopic, err)
			os.Exit(1)
		}
	}
}

func checkValidRangeOnReception(
	helperClient *mqtt.Client,
	sensorAlert config.SensorAlert,
	alertMessage string,
) pahoMqtt.MessageHandler {
	return func(mqttClient pahoMqtt.Client, message pahoMqtt.Message) {
		sensorValue := getJSONValueAsIntFromMessage(message)
		if !(sensorAlert.LowerBound <= sensorValue && sensorValue <= sensorAlert.HigherBound) {
			err := helperClient.Publish(sensorAlert.OutgoingTopic, sensorAlert.OutgoingQos, false, alertMessage)
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
