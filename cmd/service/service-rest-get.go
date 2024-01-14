package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/influxdata/influxdb-client-go/v2"
	"imt-atlantique.project.group.fr/meteo-airport/internal/logutil"
	"net/http"
)

type MeasurementHandler struct {
}

func HomeHandler(writer http.ResponseWriter, request *http.Request) {
	logutil.Info("bla blo")
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

func MeasurementInterval(w http.ResponseWriter, r *http.Request) {
	logutil.Info("This is the measurement interval handler")
	id := mux.Vars(r)["type"]
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	client := influxdb2.NewClient("http://localhost:8086",
		"upvBeRD7IGz2JkRYkF16F4PK7g-uciplnKnwMnLnFqk_5AAoT-dcUz_fWoeL0f6iy3enhBS-N0tLhwfZ0ILZiA==")

	fluxQuery := fmt.Sprintf(`from(bucket: "metrics")
	|> range(start: %s , stop: %s)
	|> filter(fn: (r) => r._measurement == "%s" )`, start, end, id)
	logutil.Info("fluxQuery: %s", fluxQuery)

	// get QueryTableResult
	result, err := client.QueryAPI("meteo-airport").Query(context.Background(), fluxQuery)
	if err != nil {
		logutil.Info("Error: %s", err)
		panic(err)
	}

	// Iterate over query response

	data := make([]map[string]interface{}, 0)
	for result.Next() {
		// Notice when group key has changed
		if result.TableChanged() {
			logutil.Info("table: %s", result.TableMetadata().String())
		}
		// Access data
		logutil.Info("value: %v", result.Record().Value())
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

func AverageMeasurementOfADay(writer http.ResponseWriter, r *http.Request) {
	logutil.Info("This is the average measurement of a day handler")
	// TODO: get measurement types from database

	day := r.URL.Query().Get("day")
	logutil.Info("day: %s", day)
	client := influxdb2.NewClient("http://localhost:8086",
		"upvBeRD7IGz2JkRYkF16F4PK7g-uciplnKnwMnLnFqk_5AAoT-dcUz_fWoeL0f6iy3enhBS-N0tLhwfZ0ILZiA==")

	fluxQuery := fmt.Sprintf(`from(bucket: "metrics")
	|> mean(column: "value")
	|> range(start: %s , stop: %s)
	|> group(columns: ["_measurement"])
	`)
	logutil.Info("fluxQuery: %s", fluxQuery)

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
			logutil.Info("table: %s", result.TableMetadata().String())
		}
		// Access data
		logutil.Info("value: %v", result.Record().Value())
		data = append(data, map[string]interface{}{
			"value": result.Record().Value(),
			"time":  result.Record().Time().String(),
			"unit":  result.Record().Field()})
	}

	jsonData, err := json.Marshal(map[string]interface{}{
		"day":  day,
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

func MeasurementOfADayByType(writer http.ResponseWriter, request *http.Request) {

}

func main() {
	router := mux.NewRouter()
	// TODO: add redirect to HomeHandler
	logutil.Info("Connected to the server on port 8081 !")
	router.HandleFunc("/api/v1/measurements", HomeHandler)
	//with query parameters
	router.HandleFunc("/api/v1/measurements/interval/{type}/", MeasurementInterval)
	router.HandleFunc("/api/v1//measurements/day/", AverageMeasurementOfADay)
	router.HandleFunc("/api/v1//measurements/byType/{type}", MeasurementOfADayByType)

	err := http.ListenAndServe(":8081", router)
	if err != nil {
		panic(err)
	}
}
