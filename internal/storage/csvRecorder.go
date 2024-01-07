package storage

import (
	"encoding/csv"
	"imt-atlantique.project.group.fr/meteo-airport/internal/logutil"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
	"os"
	"path/filepath"
	"sync"
)

type CSVRecorder struct {
	mu       sync.Mutex
	file     *os.File
	writer   *csv.Writer
	Settings CSVSettings
}

type CSVSettings struct {
	// Path is the path to the CSV file
	PathDirectory string
	Separator     rune
	TimeFormat    string
}

// NewCSVRecorder creates a new CSVRecorder with the given settings
func NewCSVRecorder(filename string, settings CSVSettings) (*CSVRecorder, error) {
	file, err := os.OpenFile(filepath.Join(settings.PathDirectory, filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logutil.Error("Failed to open file: %v", err)
		return nil, err
	}

	writer := csv.NewWriter(file)

	// Write field names in the first line if we are creating a new file
	if info, err := file.Stat(); err == nil && info.Size() == 0 {
		if err := writer.Write([]string{sensor.MeasurementFieldNames(settings.Separator)}); err != nil {
			logutil.Error("Failed to write field names: %v", err)
			return nil, err
		}

		writer.Flush()
	}

	return &CSVRecorder{
		mu:       sync.Mutex{},
		file:     file,
		writer:   writer,
		Settings: settings,
	}, nil
}

// Record stores a measurement
func (r *CSVRecorder) Record(m *sensor.Measurement) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	record := m.ToCSV(r.Settings.Separator, r.Settings.TimeFormat)

	if err := r.writer.Write([]string{record}); err != nil {
		logutil.Error("Failed to write record: %v", err)
		return err
	}

	r.writer.Flush()

	return r.writer.Error()
}

// Close closes the recorder
func (r *CSVRecorder) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.writer != nil {
		r.writer.Flush()

		if err := r.writer.Error(); err != nil {
			logutil.Error("Failed to flush writer: %v", err)
			return err
		}
	}

	if err := r.file.Close(); err != nil {
		logutil.Error("Failed to close file: %v", err)
		return err
	}

	return nil
}
