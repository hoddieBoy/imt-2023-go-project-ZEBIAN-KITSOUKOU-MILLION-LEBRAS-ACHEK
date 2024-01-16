package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	sensor "imt-atlantique.project.group.fr/meteo-airport/internal/sensor/sensor"
)



func main() {
	// Simulating humidity readings every second
	var invoiceIndex int64 = 0

	sensor.InitializeSensor()

	sensor_id = 1
	sensor.GenerateData(sensor_id, "CDG", sensor.Humidity, 30, "%", time.Now())
	if err != nil {panic(err)}

	

	for {
		humidity := readHumidity()
		invoiceIndex++
		sensor.ChangeValueMeasurement(temperature)
		fmt.Printf("Humidity: reading on\n")
		time.Sleep(4 * time.Second)
	}
}

func readHumidity() float64 {
	// Simulating humidity between 40 and 60%
	return 40 + rand.Float64() * (60 - 40)
}
