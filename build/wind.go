package build


import (
	"math/rand"
	"time"
	"encoding/json"
	"imt-2023-go-project-ZEBIAN-KITSOUKOU-MILLION-LEBRAS-ACHEK/internal"
)

func build(){
	mqttClient := mqtt.makeDefaultClient()
	actualWind := 40
	min := 10
	max := 120
	now := time.Now()

	for {
		actualWind := rand.Intn(max - min) + min

		data := map[string]interface{}{
			"idCaptor" : 2,
			"idAirport" : 1,
			"mesure" : "Wind speed",
			"value" : actualWind,
			"timestamp" : now.Format(2006-01-02) + now.Format(15-04-05),
		}
		jsonData, err := json.Marshal(data)
		
		if err != nil {
			fmt.Printf("could not marshal json: %s\n", err)
		} else {
			Publish(mqttClient, jsonData, "wind")
		}

		time.Sleep(30)
	}
}