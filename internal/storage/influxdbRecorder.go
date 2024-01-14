package storage

import (
	"context"
	"github.com/influxdata/influxdb-client-go/v2"
	"imt-atlantique.project.group.fr/meteo-airport/internal/config_helper"
	"imt-atlantique.project.group.fr/meteo-airport/internal/logutil"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
	"strconv"
	"sync"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"imt-atlantique.project.group.fr/meteo-airport/internal/log"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
)

type InfluxDBRecorder struct {
	mu     sync.Mutex
	client influxdb2.Client
	bucket string
	org    string
}

func NewInfluxDBRecorder(settings config_helper.InfluxDBSettings) (*InfluxDBRecorder, error) {
	client := influxdb2.NewClient(settings.URL, settings.Token)

	return &InfluxDBRecorder{
		mu:     sync.Mutex{},
		client: client,
		bucket: settings.Bucket,
		org:    settings.Organization,
	}, nil
}

func (r *InfluxDBRecorder) RecordOnContext(ctx context.Context, m *sensor.Measurement) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	writeAPI := r.client.WriteAPIBlocking(r.org, r.bucket)

	p := influxdb2.NewPointWithMeasurement(string(m.Type)).
		AddTag("airport", m.AirportID).
		AddTag("sensor", strconv.FormatInt(m.SensorID, 10)).
		AddTag("unit", m.Unit).
		AddField("value", m.Value).
		SetTime(m.Timestamp)

	if err := writeAPI.WritePoint(ctx, p); err != nil {
		log.Error("Failed to write point on InfluxDB: %v", err)
		return err
	}

	return nil
}

func (r *InfluxDBRecorder) Record(m *sensor.Measurement) error {
	return r.RecordOnContext(context.Background(), m)
}

func (r *InfluxDBRecorder) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.client.Close()

	return nil
}
