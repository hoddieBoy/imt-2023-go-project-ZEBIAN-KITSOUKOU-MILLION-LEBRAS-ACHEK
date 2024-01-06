package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	internal "imt-atlantique.project.group.fr/meteo-airport/internal"
)

func main() {
	client := internal.MakeDefaultClient()
	// MQTT broker configuration

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
