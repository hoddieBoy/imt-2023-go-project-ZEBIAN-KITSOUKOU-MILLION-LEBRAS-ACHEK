package main

import (
	"fmt"
	"imt-atlantique.project.group.fr/meteo-airport/internal/logutil"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
	"math/rand"
	"time"

	mqtt_helper "imt-atlantique.project.group.fr/meteo-airport/internal/mqtt_helper"
)

func windDataGeneration(actualWind float64, min float64, max float64) float64 {
	actualWind = actualWind + (rand.Float64()-rand.Float64())*5

	if actualWind < min {
		actualWind = min
	}
	if actualWind > max {
		actualWind = max
	}

	return actualWind
}

func publishData(actualWind float64, client *mqtt_helper.MQTTClient) {

	data := sensor.Measurement{
		SensorID:  2,
		AirportID: "CDG",
		Type:      "Wind speed",
		Value:     actualWind,
		Unit:      "Km/h",
		Timestamp: time.Now(),
	}

	err := data.PublishOnMQTT(2, false, client)
	if err != nil {
		logutil.Error(fmt.Sprintf("Failed to publish data to client: %v", err))
	}
}

func main() {
	if config, err := mqtt_helper.RetrieveMQTTPropertiesFromYaml("./config/hiveClientConfig.yaml"); err != nil {
		panic(err)
	} else {
		client := mqtt_helper.NewClient(config, "clientId")

		err := client.Connect()
		if err != nil {
			logutil.Error(fmt.Sprintf("Cannot connect to client: %v", err))
		}

		actualWind := 40.0
		minimalValue := 10.0
		maximalValue := 120.0

		for {

			actualWind := windDataGeneration(actualWind, minimalValue, maximalValue)
			publishData(actualWind, client)

			time.Sleep(5 * time.Second)
		}
	}
}
