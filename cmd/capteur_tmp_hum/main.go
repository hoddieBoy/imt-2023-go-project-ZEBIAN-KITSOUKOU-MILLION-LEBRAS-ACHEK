package main

import (
	"fmt"
	"math/rand"
	"time"
	
	mqtt "github.com/eclipse/paho.mqtt.golang"
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
	// MQTT broker configuration
	brokerURI := "mqtt://localhost:1883" //maybe tcp://localhost:1883
	clientID := "jam_client"
	// username := "emqx"
	// password := ""

	client := connect(brokerURI, clientID)
	// Send fake temperature and humidity data every 3 seconds for temperature and 5 seconds for humidity
	for {
		// Generate random temperature and humidity values
		temperature := 20 + rand.Float64() * (30 - 20)
		humidity := 40 + rand.Float64() * (60 - 40)

		// Publish temperature data to MQTT topic
		token := client.Publish("capteur/T", 0, false, fmt.Sprintf("%.2f", temperature))
		token.Wait()

		// Publish humidity data to MQTT topic
		token = client.Publish("capteur/H", 0, false, fmt.Sprintf("%.2f", humidity))
		token.Wait()

		// Sleep for 3 seconds for temperature and 5 seconds for humidity
		time.Sleep(3 * time.Second)
		time.Sleep(5 * time.Second)
	}
}

