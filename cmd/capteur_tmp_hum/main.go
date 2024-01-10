package main

import (
	"fmt"
	"math/rand"
	"time"
	
	mqtt_helper "imt-atlantique.project.group.fr/meteo-airport/internal/mqtt_helper"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	sensor "imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
)


func createClientOptions(brokerURI string, clientId string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURI)
	// opts.SetUsername(username)
	// opts.SetPassword(password)
	opts.SetClientID(clientId)
	return opts
}

func connect(brokerURI string, clientId string) mqtt.Client {
	opts := createClientOptions(brokerURI, clientId)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		panic(err)
	}
	return client
}

func main() {
	config, err := mqtt_helper.RetrieveMQTTPropertiesFromYaml("./config/hiveClientConfig.yaml")

	if err != nil {panic(err)}

	var invoiceIndex int64 = 0
	measurementT := sensor.Measurement{
		SensorID:  1,
		AirportID: "CDG",
		Type:      sensor.Temperature,
		Value:     20,
		Unit:      "°C",
		Timestamp: time.Now(),
	}

	measurementH := sensor.Measurement{
		SensorID:  invoiceIndex,
		AirportID: "CDG",
		Type:      sensor.Humidity,
		Value:     30,
		Unit:      "%",
		Timestamp: time.Now(),
	}


	clientID := "jam_client"

	client := mqtt_helper.NewClient(config, clientID)
	client.Connect()
	// username := "emqx"
	// password := ""

	// Send fake temperature and humidity data every 3 seconds for temperature and 5 seconds for humidity
	for {
		// Generate random temperature and humidity values
		temperature := 20 + rand.Float64() * (30 - 20)
		humidity := 40 + rand.Float64() * (60 - 40)

	invoiceIndex++

	measurementT.SensorID = invoiceIndex
	measurementT.Timestamp = time.Now()
	measurementT.Value = temperature

	measurementH.SensorID = invoiceIndex
	measurementH.Timestamp = time.Now()
	measurementH.Value = humidity

		// Publish temperature data to MQTT topic
		// token := client.Publish("capteur/T", 2, false, fmt.Sprintf("%.2f", temperature))
		measurementT.PublishOnMQTT(2, false, client)
		// token.Wait()
		fmt.Printf("%#v",measurementT)

		// Publish humidity data to MQTT topic
		// token = client.Publish("capteur/H", 2, false, fmt.Sprintf("%.2f", humidity))
		measurementH.PublishOnMQTT(2, false, client)
		fmt.Printf("%#v",measurementH)
		// token.Wait()

		// Sleep for 3 seconds for temperature and 5 seconds for humidity
		time.Sleep(4 * time.Second)
	}
}

