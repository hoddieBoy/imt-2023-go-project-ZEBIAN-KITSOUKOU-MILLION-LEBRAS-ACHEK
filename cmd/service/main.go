// File: main.go
// Description: main.go is the entry point of the API.
// It is responsible for handling the routes and launching the server

// @title Meteo Airport API
// @description This is the API for the Meteo Airport project.
// @version 1
// @host localhost:8082
// @BasePath /api/v1/measurements
// @schemes http

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	_ "imt-atlantique.project.group.fr/meteo-airport/cmd/service/docs"

	"github.com/gorilla/mux"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	httpSwagger "github.com/swaggo/http-swagger"
	"imt-atlantique.project.group.fr/meteo-airport/internal/log"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
)

var client = influxdb2.NewClient("http://localhost:8086",
	"upvBeRD7IGz2JkRYkF16F4PK7g-uciplnKnwMnLnFqk_5AAoT-dcUz_fWoeL0f6iy3enhBS-N0tLhwfZ0ILZiA==")
var newURL = "http://localhost:8082/api/v1/measurements"

func RedirectHomeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, newURL, http.StatusSeeOther)
}

func HomeHandler(writer http.ResponseWriter, _ *http.Request) {
	jsonData, err := json.Marshal(map[string]interface{}{
		"message": "Hello welcome to our API",
	},
	)

	handleErr(err)

	setResponseHeaders(writer)
	_, err = writer.Write(jsonData)

	handleErr(err)
}

// MeasurementIntervalHandler gets measurements in a specified interval.
// @Summary Get measurements in a specific time interval
// @Description Get measurements for a specified type within a time range.
// @ID measurement-interval
// @Accept json
// @Produce json
// @Param type path string true "Measurement type"
// @Param start query string true "Start date in the format -Hh"
// @Param end query string true "End date in the format -Hh"
// @Success 200 {object} map[string]interface{}
// @Router /interval/{type}/ [get]
func MeasurementIntervalHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["type"]
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	fluxQuery := fmt.Sprintf(`from(bucket: "metrics")
	|> range(start: %s , stop: %s)
	|> filter(fn: (r) => r._measurement == "%s" )`, start, end, id)

	result, err := client.QueryAPI("meteo-airport").Query(context.Background(), fluxQuery)
	handleErr(err)

	data := make([]map[string]interface{}, 0)

	for result.Next() {
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
	handleErr(err)

	setResponseHeaders(w)
	_, err = w.Write(jsonData)

	handleErr(err)
}

// AvgMeasurementInADayHandler gets the average measurement for specified types on a given date.
// @Summary Get average measurement in a day
// @Description Get the average measurement for specified types on a given date
// @ID avg-measurement-in-a-day
// @Accept json
// @Produce json
// @Param date query string true "Date in the format YYYY-MM-DD"
// @Param types query string false "Comma-separated list of measurement types
// (e.g., temperature,humidity,pressure,windSpeed)"
// @Success 200 {object} map[string]interface{}
// @Router /mean/ [get]
func AvgMeasurementInADayHandler(writer http.ResponseWriter, r *http.Request) {
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

	result, err := client.QueryAPI("meteo-airport").Query(context.Background(), fluxQuery)
	handleErr(err)

	data := processData(*result)
	jsonData, err := json.Marshal(map[string]interface{}{
		"date": date,
		"data": data,
	})

	handleErr(err)

	setResponseHeaders(writer)
	_, err = writer.Write(jsonData)

	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func setResponseHeaders(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}

func processData(result api.QueryTableResult) []map[string]interface{} {
	data := make([]map[string]interface{}, 0)

	for result.Next() {
		var measurementType string

		switch result.Record().Measurement() {
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

	return data
}

func main() {
	router := mux.NewRouter()

	log.Info("Connected to the server on port 8082 !")
	router.HandleFunc("/", RedirectHomeHandler)
	router.HandleFunc("/api", RedirectHomeHandler)
	router.HandleFunc("/api/v1", RedirectHomeHandler)
	router.HandleFunc("/api/v1/measurements", HomeHandler).Methods("GET")
	router.HandleFunc("/api/v1/measurements/interval/{type}/", MeasurementIntervalHandler).Methods("GET")
	router.HandleFunc("/api/v1/measurements/mean/", AvgMeasurementInADayHandler).Methods("GET")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	server := &http.Server{
		Addr:              ":8082",
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}
