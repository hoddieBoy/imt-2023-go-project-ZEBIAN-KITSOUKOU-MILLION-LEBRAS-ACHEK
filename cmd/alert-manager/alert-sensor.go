package main

import (
	"fmt"

	pahoMqtt "github.com/eclipse/paho.mqtt.golang"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
)

func main() {
	config, err := mqtt.RetrieveMQTTPropertiesFromYaml("./config/hiveClientConfig.yaml")
	if err != nil {
		panic(err)
	}

	client := mqtt.NewClient(config, "anotherClientId")
	if err := client.Connect(); err != nil {
		panic(err)
	}

	defer client.Disconnect()

	handleAlertListening(client)

	select {}
}

func handleAlertListening(client *mqtt.Client) {
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

func checkValidRangeOnReception(
	helperClient *mqtt.Client,
	sensorAlert mqtt.SensorAlertType,
	alertMessage string,
) pahoMqtt.MessageHandler {
	return func(mqttClient pahoMqtt.Client, message pahoMqtt.Message) {
		sensorValue := getJSONValueAsIntFromMessage(message)
		if !(sensorAlert.LowerBound <= sensorValue && sensorValue <= sensorAlert.HigherBound) {
			err := helperClient.Publish(sensorAlert.EndPoint, 1, false, alertMessage)
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
