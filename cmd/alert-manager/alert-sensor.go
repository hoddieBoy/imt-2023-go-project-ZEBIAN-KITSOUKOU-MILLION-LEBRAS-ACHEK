package main

import (
	"fmt"

	"os"

	pahoMqtt "github.com/eclipse/paho.mqtt.golang"
	"imt-atlantique.project.group.fr/meteo-airport/internal/config"
	"imt-atlantique.project.group.fr/meteo-airport/internal/log"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		log.Warn("No config file specified, using default path: config/config.yaml")
	} else {
		config.SetDefaultConfigFileName(args[0])
	}

	alertConfig, err := config.LoadDefaultAlertConfig()

	if err != nil {
		panic(err)
	}

	clients := make(map[string]*mqtt.Client)

	for sensorType, alert := range alertConfig.SensorsAlert {
		clients[sensorType] = mqtt.NewClient(&alertConfig.Broker, alert.ClientID)
		err := clients[sensorType].Connect()
		if err != nil {
			log.Error("Error connecting to broker: %v", err)
			os.Exit(1)
		}
	}

	defer func() {
		for _, client := range clients {
			client.Disconnect()
		}
	}()

	handleAlertListening(clients, alertConfig.SensorsAlert)

	select {}
}

func handleAlertListening(clients map[string]*mqtt.Client, alerts map[string]config.SensorAlert) {
	for sensorType, alert := range alerts {
		err := clients[sensorType].Subscribe(
			alert.IncomingTopic,
			alert.IncomingQos,
			checkValidRangeOnReception(clients[sensorType], alert,
				"Alert: "+sensorType+" value is out of range ["+fmt.Sprintf("%f", alert.LowerBound)+", "+fmt.Sprintf("%f", alert.HigherBound)+"]"),
		)
		if err != nil {
			log.Error("An error occured while subscribing to topic %s: %v", alert.IncomingTopic, err)
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

func getJSONValueAsIntFromMessage(message pahoMqtt.Message) float64 {
	measurement, err := sensor.FromJSON(message.Payload())
	if err != nil {
		log.Warn("Unable to retrieve value from message: %v", err)
		return 0
	}
	return measurement.Value
}
