package storage

import (
	"encoding/csv"
	"fmt"
	"imt-atlantique.project.group.fr/meteo-airport/internal/config_helper"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"imt-atlantique.project.group.fr/meteo-airport/internal/log"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
)

type CSVRecorder struct {
	mu       sync.Mutex
	file     *os.File
	writer   *csv.Writer
	Settings config_helper.CSVSettings
}

func NewCSVRecorder(settings config_helper.CSVSettings) (*CSVRecorder, error) {
	return &CSVRecorder{
		mu:       sync.Mutex{},
		Settings: settings,
	}, nil
}

func (r *CSVRecorder) setWriter(filename string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Open or create file
	file, err := os.OpenFile(filepath.Join(r.Settings.PathDirectory, filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	r.file = file
	r.writer = csv.NewWriter(file)

	// Write header if file is empty
	if fileInfo, err := file.Stat(); err == nil && fileInfo.Size() == 0 {
		if err := r.writer.Write([]string{sensor.MeasurementFieldNames(r.Settings.Separator)}); err != nil {
			return fmt.Errorf("failed to write header: %w", err)
		}
		r.writer.Flush()
	}

	return nil
}

func (r *CSVRecorder) Record(m *sensor.Measurement) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	record := m.ToCSV(r.Settings.Separator, r.Settings.TimeFormat)
	filename := "airport_" + m.AirportID + "_sensor_" + strconv.FormatInt(m.SensorID, 10) + "_" + string(m.Type) + "_" + m.Timestamp.Format("2006-01-02") + ".csv"

	if r.writer == nil || r.file.Name() != filepath.Join(r.Settings.PathDirectory, filename) {
		if err := r.setWriter(filename); err != nil {
			return err
		}
	}

	if err := r.writer.Write([]string{record}); err != nil {
		return fmt.Errorf("failed to write record: %w", err)
	}

	r.writer.Flush()

	return nil
}

func (r *CSVRecorder) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.writer != nil {
		r.writer.Flush()

		if err := r.writer.Error(); err != nil {
			return fmt.Errorf("failed to flush writer: %w", err)
		}
	}
	if err := r.file.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	return nil
}
