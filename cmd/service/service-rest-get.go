package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/influxdata/influxdb-client-go/v2"
	"imt-atlantique.project.group.fr/meteo-airport/internal/log"
)

type MeasurementHandler struct {
}

func HomeHandler(writer http.ResponseWriter, request *http.Request) {
	log.Info("This is the home handler")
	jsonData, err := json.Marshal(map[string]interface{}{

		"message": "Hello welcome to our API",
	},
	)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		panic(err)
	}
}

func AverageMeasurementInADayHandler(writer http.ResponseWriter, r *http.Request) {
	log.Info("This is the average measurement in a day handler")

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

	day := r.URL.Query().Get("day")

	client := influxdb2.NewClient("http://localhost:8086",
		"upvBeRD7IGz2JkRYkF16F4PK7g-uciplnKnwMnLnFqk_5AAoT-dcUz_fWoeL0f6iy3enhBS-N0tLhwfZ0ILZiA==")

	var typeF []string
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
	`, day, day+"T23:59:59Z", strings.Join(typeF, " or "))
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
		if result.Record().Measurement() == "temperature" {
			measurementType = "temperature"
		} else if result.Record().Measurement() == "humidity" {
			measurementType = "humidity"
		} else if result.Record().Measurement() == "pressure" {
			measurementType = "pressure"
		} else if result.Record().Measurement() == "windSpeed" {
			measurementType = "windSpeed"
		}
		data = append(data, map[string]interface{}{
			"type":  measurementType,
			"value": result.Record().Value(),
			"time":  result.Record().Time().String(),
			"unit":  result.Record().Field(),
		})
	}

	jsonData, err := json.Marshal(map[string]interface{}{
		"data": data,
	})

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(jsonData)
	if err != nil {
		panic(err)
	}

}

func main() {
	router := mux.NewRouter()
	// TODO: add redirect to HomeHandler
	log.Info("Connected to the server on port 8081 !")
	router.HandleFunc("/api/v1/measurements", HomeHandler)
	router.HandleFunc("/api/v1/measurements/interval/{type}/", MeasurementIntervalHandler)
	router.HandleFunc("/api/v1/measurements/mean/", AverageMeasurementInADayHandler)

	err := http.ListenAndServe(":8081", router)
	if err != nil {
		panic(err)
	}
}
