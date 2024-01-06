package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	internal "imt-atlantique.project.group.fr/meteo-airport/internal"
)

/*func createClientOptions(brokerURI string, clientId string) *mqtt.ClientOptions {
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
}*/

func main() {
	client := internal.MakeDefaultClient()
	// MQTT broker configuration
	/*brokerURI := "mqtt://411d6c045163486b846c891f3910e83f.s2.eu.hivemq.cloud:1883" //maybe tcp://localhost:1883
	clientID := "greg_client"

	client := connect(brokerURI, clientID)*/
	actualWind := 40
	min := 10
	max := 120

	type wind struct {
		idCaptor  int
		idAirport string
		mesure    string
		value     int
		timestamp time.Time
	}

	for {
		actualWind = rand.Intn(max-min) + min

		data := wind{
			idCaptor:  2,
			idAirport: "CDG",
			mesure:    "Wind speed",
			value:     actualWind,
			timestamp: time.Now(),
		}
		jsonData, err := json.Marshal(data)

		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
		} else {
			internal.Publish(client, fmt.Sprintf("%.2f", jsonData), "capteur/W")
			fmt.Printf("data send")
		}

		time.Sleep(5)
	}
}
