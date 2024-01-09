package storage

import "imt-atlantique.project.group.fr/meteo-airport/internal/sensor"

type Recorder interface {
	Record(m *sensor.Measurement) error

	Close() error
}
