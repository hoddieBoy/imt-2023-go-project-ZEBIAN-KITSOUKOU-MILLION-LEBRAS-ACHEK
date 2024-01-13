package main

import (
	"fmt"
	pahoMqtt "github.com/eclipse/paho.mqtt.golang"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
)

func main() {
	if config, err := mqtt.RetrieveMQTTPropertiesFromYaml("./config/hiveClientConfig.yaml"); err != nil {
		panic(err)
	} else {
		client := mqtt.NewClient(config, "anotherClientId")
		if err := client.Connect(); err != nil {
			panic(err)
		}
		defer client.Disconnect()

		handleAlertListening(client)

		select {}
	}
}

func handleAlertListening(client *mqtt.MQTTClient) {
	rootConfig, err := mqtt.RetrieveMQTTRootFromYaml()
	if err != nil {
		panic(err)
	}

	err = client.Subscribe(rootConfig.Root.Sensor.Humidity,
		1,
		checkValidRangeOnReception(client,
			rootConfig.Root.Alert.Humidity,
			"Alert, Humidity sensor out of range"))
	if err != nil {
		panic(err)
	}

	err = client.Subscribe(rootConfig.Root.Sensor.Temperature,
		1,
		checkValidRangeOnReception(client,
			rootConfig.Root.Alert.Temperature,
			"Alert, Temperature sensor out of range"))
	if err != nil {
		panic(err)
	}

	err = client.Subscribe(rootConfig.Root.Sensor.Pressure,
		1,
		checkValidRangeOnReception(client,
			rootConfig.Root.Alert.Pressure,
			"Alert, pressure sensor out of range"))
	if err != nil {
		panic(err)
	}
}

func checkValidRangeOnReception(helperClient *mqtt.MQTTClient, sensorAlert mqtt.SensorAlertType, alertMessage string) pahoMqtt.MessageHandler {
	return func(mqttClient pahoMqtt.Client, message pahoMqtt.Message) {
		sensorValue := getJsonValueAsIntFromMessage(message)
		if !(sensorAlert.LowerBound <= sensorValue && sensorValue <= sensorAlert.HigherBound) {
			err := helperClient.Publish(sensorAlert.EndPoint, 1, false, alertMessage)
			if err != nil {
				panic(err)
			}
		}
	}
}

func getJsonValueAsIntFromMessage(message pahoMqtt.Message) int {
	// TODO
	fmt.Println(string(message.Payload()))
	return 50
}
