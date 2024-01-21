package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"imt-atlantique.project.group.fr/meteo-airport/internal/logutil"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt_helper"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
	_ "imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
)

func main() {
	if config, err := mqtt_helper.RetrieveMQTTPropertiesFromYaml("./config/hiveClientConfig.yaml"); err != nil {
		panic(err)
	} else {
		client := mqtt_helper.NewClient(config, "anotherClientId")
		if err := client.Connect(); err != nil {
			panic(err)
		}
		defer client.Disconnect()

		handleAlertListening(client)

		for true {

		}
	}
}

func handleAlertListening(client *mqtt_helper.MQTTClient) {
	rootConfig, err := mqtt_helper.RetrieveMQTTRootFromYaml()
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

func checkValidRangeOnReception(helperClient *mqtt_helper.MQTTClient, sensorAlert mqtt_helper.SensorAlertType, alertMessage string) func(client mqtt.Client, message mqtt.Message) {
	return func(mqttClient mqtt.Client, message mqtt.Message) {
		sensorValue := getJsonValueAsIntFromMessage(message)
		if !(float64(sensorAlert.LowerBound) <= sensorValue && sensorValue <= float64(sensorAlert.HigherBound)) {
			err := helperClient.Publish(sensorAlert.EndPoint, 1, false, alertMessage)
			if err != nil {
				panic(err)
			}
		}
	}
}

func getJsonValueAsIntFromMessage(message mqtt.Message) float64 {
	measurement, err := sensor.FromJSON(message.Payload())
	if err != nil {
		logutil.Warn("Error unmarshalling measurement from JSON: %v", err)
		return 0
	}
	return measurement.Value
}
