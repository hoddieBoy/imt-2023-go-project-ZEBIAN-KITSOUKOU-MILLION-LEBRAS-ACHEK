package storage

import "imt-atlantique.project.group.fr/meteo-airport/internal/sensor"

type Recorder interface {
	// Record stores a measurement
	Record(m *sensor.Measurement) error

	// Close closes the recorder
	Close() error
}
