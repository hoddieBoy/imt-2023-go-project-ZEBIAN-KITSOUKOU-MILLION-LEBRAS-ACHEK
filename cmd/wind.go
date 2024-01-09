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
	//internal.SubscribeWithQos_1(client, "capteur/W")

	actualWind := 40
	min := 10
	max := 120

	type wind struct {
		idCaptor  int
		idAirport string
		mesure    string
		value     int
		timestamp "2006-01-02-15-04-05"
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
			fmt.Printf("%d", jsonData)
		}

		time.Sleep(5000)
	}
}
