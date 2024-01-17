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
			checkValidRangeOnReception(clients[sensorType], alert, sensorType),
		)

		if err != nil {
			log.Error("An error occurred while subscribing to topic %s: %v", alert.IncomingTopic, err)
			os.Exit(1)
		}
	}
}

func checkValidRangeOnReception(
	helperClient *mqtt.Client,
	sensorAlert config.SensorAlert,
	sensorType string,
) pahoMqtt.MessageHandler {
	return func(mqttClient pahoMqtt.Client, message pahoMqtt.Message) {
		sensorValue := getJSONValueAsIntFromMessage(message)
		alertMessage := "Alert: " + sensorType + " value is "

		if sensorValue < sensorAlert.LowerBound {
			alertMessage += "lower than " + fmt.Sprintf("%f", sensorAlert.LowerBound)
		} else {
			alertMessage += "higher than " + fmt.Sprintf("%f", sensorAlert.HigherBound)
		}

		err := helperClient.Publish(sensorAlert.OutgoingTopic, sensorAlert.OutgoingQos, false, alertMessage)
		if err != nil {
			log.Warn("Error publishing alert message: %v", err)
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
