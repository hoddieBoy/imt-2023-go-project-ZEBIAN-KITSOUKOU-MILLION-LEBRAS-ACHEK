package main

import (
	"fmt"
	"math/rand"
	"time"

	"imt-atlantique.project.group.fr/meteo-airport/internal/log"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
)



func main() {
	// Simulating humidity readings every second
	var invoiceIndex int64 = 0
	var sensor_id int64 = 1
	var humiditySetup float64 = 40.0

	sensor := sensor.Sensor{}
	err := sensor.InitializeSensor(sensor_id, "CDG", "humidity", 30, "%", time.Now())
//	sensor.GenerateData(sensor_id, "CDG", sensor.Humidity, 30, "%", time.Now())
	if err != nil {panic(err)}

	for {
		humidity := readHumidity(humiditySetup)
		invoiceIndex++
		sensor.ChangeValueMeasurement(humidity)
		err := sensor.PublishData()
		if err != nil {
			log.Error(fmt.Sprintf("Failed to publish data to client: %v", err))
		}
		fmt.Printf("Humidity: %f\n", humidity)
		time.Sleep(4 * time.Second)
	}
}

func readHumidity(currentHumidity float64) float64 {
	// Simulating humidity between 40 and 60%

	simulatedHumidity := currentHumidity

	if (simulatedHumidity < 40) {
		simulatedHumidity = 40
	}
	if (simulatedHumidity > 60) {
		simulatedHumidity = 60
	}

	simulatedHumidity = simulatedHumidity + rand.Float64() * 2

	return simulatedHumidity
}
