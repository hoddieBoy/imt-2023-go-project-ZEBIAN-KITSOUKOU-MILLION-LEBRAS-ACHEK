package storage

import (
	"fmt"
	"sync"

	pahoMqtt "github.com/eclipse/paho.mqtt.golang"
	"imt-atlantique.project.group.fr/meteo-airport/internal/log"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
)

// Manager struct handles subscriptions and manages recorders
type Manager struct {
	recordersMutex sync.RWMutex
	recorders      map[sensor.MeasurementType]map[Recorder]Topic
	mqttClient     *mqtt.Client
	closeChan      chan struct{}
	wg             sync.WaitGroup
}

type Topic struct {
	Qos  byte
	name string
}

func NewManager(mqttClient *mqtt.Client) *Manager {
	return &Manager{
		recorders:  make(map[sensor.MeasurementType]map[Recorder]Topic),
		mqttClient: mqttClient,
		closeChan:  make(chan struct{}),
		wg:         sync.WaitGroup{},
	}
}

func (s *Manager) AddRecorder(sensorType sensor.MeasurementType, topic string, qos byte, recorder Recorder) error {
	s.recordersMutex.Lock()
	if _, ok := s.recorders[sensorType]; !ok {
		s.recorders[sensorType] = make(map[Recorder]Topic)
	}

	s.recorders[sensorType][recorder] = Topic{Qos: qos, name: topic}

	s.recordersMutex.Unlock()

	return nil
}

func (s *Manager) subscribeToSensor(sensorType sensor.MeasurementType, topic string, qos byte) error {
	return s.mqttClient.Subscribe(topic, qos, func(client pahoMqtt.Client, message pahoMqtt.Message) {
		measurement, err := sensor.FromJSON(message.Payload())
		if err != nil {
			log.Warn("Error unmarshalling measurement from JSON: %v", err)
			return
		}

		s.recordersMutex.RLock()
		defer s.recordersMutex.RUnlock()

		for recorder := range s.recorders[sensorType] {
			s.wg.Add(1)
			go func(rec Recorder, meas *sensor.Measurement) {
				if err := rec.Record(meas); err != nil {
					log.Warn("Error recording measurement of type %s with recorder %v: %v", sensorType, rec, err)
				}
				s.wg.Done()
			}(recorder, measurement)
		}
	})
}

func (s *Manager) Close() error {
	close(s.closeChan)

	s.wg.Wait()

	var closeErrors []error

	s.recordersMutex.RLock()
	defer s.recordersMutex.RUnlock()

	for _, recorders := range s.recorders {
		for recorder := range recorders {
			if err := recorder.Close(); err != nil {
				closeErrors = append(closeErrors, fmt.Errorf("failed to close recorder: %v", err))
			}
		}
	}

	if len(closeErrors) > 0 {
		return fmt.Errorf("encountered errors while closing recorders: %v", closeErrors)
	}

	return nil
}

func (s *Manager) Start() error {
	s.recordersMutex.RLock()
	defer s.recordersMutex.RUnlock()

	for sensorType, recorders := range s.recorders {
		for _, topic := range recorders {
			if err := s.subscribeToSensor(sensorType, topic.name, topic.Qos); err != nil {
				return err
			}
		}
	}

	return nil
}
