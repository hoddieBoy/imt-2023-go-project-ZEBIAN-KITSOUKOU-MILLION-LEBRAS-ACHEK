package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"imt-atlantique.project.group.fr/meteo-airport/internal/log"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
)

var newURL = "http://localhost:8081/api/v1/measurements"
var content = "Content-Type"
var application = "application/json"

func RedirectHomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("This is the redirect home handler")
	http.Redirect(w, r, newURL, http.StatusSeeOther)
}

func HomeHandler(writer http.ResponseWriter, _ *http.Request) {
	log.Info("This is the home handler")

	jsonData, err := json.Marshal(map[string]interface{}{
		"message": "Hello welcome to our API",
	},
	)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set(content, application)
	_, err = writer.Write(jsonData)

	if err != nil {
		return
	}
}

func MeasurementIntervalHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("This is the measurement interval handler")

	id := mux.Vars(r)["type"]
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	client := influxdb2.NewClient("http://localhost:8086",
		"upvBeRD7IGz2JkRYkF16F4PK7g-uciplnKnwMnLnFqk_5AAoT-dcUz_fWoeL0f6iy3enhBS-N0tLhwfZ0ILZiA==")

	fluxQuery := fmt.Sprintf(`from(bucket: "metrics")
	|> range(start: %s , stop: %s)
	|> filter(fn: (r) => r._measurement == "%s" )`, start, end, id)
	log.Info("fluxQuery: %s", fluxQuery)

	// get QueryTableResult
	result, err := client.QueryAPI("meteo-airport").Query(context.Background(), fluxQuery)
	if err != nil {
		log.Info("Error: %s", err)
		panic(err)
	}

	// Iterate over query response

	data := make([]map[string]interface{}, 0)

	for result.Next() {
		// Notice when group key has changed
		if result.TableChanged() {
			log.Info("table: %s", result.TableMetadata().String())
		}
		// Access data
		log.Info("value: %v", result.Record().Value())
		data = append(data, map[string]interface{}{"value": result.Record().Value(),
			"time": result.Record().Time().String(),
			"unit": result.Record().Field()})
	}

	jsonData, err := json.Marshal(map[string]interface{}{
		"type":  id,
		"start": start,
		"end":   end,
		"data":  data,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(content, application)
	_, err = w.Write(jsonData)

	if err != nil {
		panic(err)
	}
}

func AvgMeasurementInADayHandler(writer http.ResponseWriter, r *http.Request) {
	log.Info("This is the average measurement in a date handler")

	var types []string

	t := r.URL.Query().Get("types")

	if t == "" {
		types = []string{"temperature", "humidity", "pressure", "windSpeed"}
	} else {
		types = strings.Split(t, ",")
	}

	if len(types) == 0 {
		types = append(types, "temperature", "humidity", "pressure", "windSpeed")
	}

	date := r.URL.Query().Get("date")

	client := influxdb2.NewClient("http://localhost:8086",
		"upvBeRD7IGz2JkRYkF16F4PK7g-uciplnKnwMnLnFqk_5AAoT-dcUz_fWoeL0f6iy3enhBS-N0tLhwfZ0ILZiA==")

	typeF := make([]string, 0)
	for _, t := range types {
		typeF = append(typeF, fmt.Sprintf(`r._measurement == "%s"`, t))
	}

	fluxQuery := fmt.Sprintf(`from(bucket: "metrics")
	|> range(start: %s, stop: %s)
	|> filter(fn: (r) => %s)
	|> filter(fn: (r) => r["_field"] == "value")
	|> aggregateWindow(every: 24h, fn: mean, createEmpty: false)
	|> group(columns: ["_measurement"])
	|> yield(name: "mean")
	`, date, date+"T23:59:59Z", strings.Join(typeF, " or "))
	log.Info("fluxQuery: %s", fluxQuery)

	// get QueryTableResult
	result, err := client.QueryAPI("meteo-airport").Query(context.Background(), fluxQuery)
	if err != nil {
		panic(err)
	}

	// Iterate over query response
	data := make([]map[string]interface{}, 0)

	for result.Next() {
		// Notice when group key has changed
		if result.TableChanged() {
			log.Info("table: %s", result.TableMetadata().String())
		}

		// Access data
		log.Info("value: %v", result.Record())

		var measurementType string

		switch measurementType {
		case string(sensor.Temperature):
			measurementType = string(sensor.Temperature)
		case string(sensor.Humidity):
			measurementType = string(sensor.Humidity)
		case string(sensor.Pressure):
			measurementType = string(sensor.Pressure)
		case string(sensor.WindSpeed):
			measurementType = string(sensor.WindSpeed)
		}

		data = append(data, map[string]interface{}{
			"type":  measurementType,
			"value": result.Record().Value(),
			"unit":  result.Record().Field(),
		})
	}

	jsonData, err := json.Marshal(map[string]interface{}{
		"date": date,
		"data": data,
	})

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set(content, application)
	_, err = writer.Write(jsonData)

	if err != nil {
		panic(err)
	}
}

func main() {
	router := mux.NewRouter()

	log.Info("Connected to the server on port 8081 !")
	router.HandleFunc("/", RedirectHomeHandler)
	router.HandleFunc("/api", RedirectHomeHandler)
	router.HandleFunc("/api/v1", RedirectHomeHandler)
	router.HandleFunc("/api/v1/measurements", HomeHandler)
	router.HandleFunc("/api/v1/measurements/interval/{type}/", MeasurementIntervalHandler)
	router.HandleFunc("/api/v1/measurements/mean/", AvgMeasurementInADayHandler)

	server := &http.Server{
		Addr:              ":8081",
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
