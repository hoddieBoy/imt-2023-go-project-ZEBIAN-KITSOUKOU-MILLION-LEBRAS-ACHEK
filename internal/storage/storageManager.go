package storage

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"imt-atlantique.project.group.fr/meteo-airport/internal/logutil"
	"imt-atlantique.project.group.fr/meteo-airport/internal/mqtt_helper"
	"imt-atlantique.project.group.fr/meteo-airport/internal/sensor"
	"sync"
)

// Manager struct handles subscriptions and manages recorders
type Manager struct {
	recordersMutex sync.RWMutex
	recorders      map[sensor.Type]map[Recorder]byte
	mqttClient     *mqtt_helper.MQTTClient
	closeChan      chan struct{}
	wg             sync.WaitGroup
}

// NewManager creates a new Manager instance
func NewManager(mqttClient *mqtt_helper.MQTTClient) *Manager {
	return &Manager{
		recorders:  make(map[sensor.Type]map[Recorder]byte),
		mqttClient: mqttClient,
		closeChan:  make(chan struct{}),
		wg:         sync.WaitGroup{},
	}
}

// AddRecorder adds a recorder for a given sensor type
func (s *Manager) AddRecorder(sensorType sensor.Type, recorder Recorder, qos byte) {
	s.recordersMutex.Lock()
	if _, ok := s.recorders[sensorType]; !ok {
		s.recorders[sensorType] = make(map[Recorder]byte)
	}
	s.recorders[sensorType][recorder] = qos
	s.recordersMutex.Unlock()
}

// SubscribeToSensor subscribes to the MQTT topic for a given sensor type
func (s *Manager) SubscribeToSensor(sensorType sensor.Type, qos byte) error {
	return s.mqttClient.Subscribe(sensorType.GetTopic(), qos, func(client mqtt.Client, message mqtt.Message) {
		measurement, err := sensor.FromJSON(message.Payload())
		if err != nil {
			logutil.Warn("Error unmarshalling measurement from JSON: %v", err)
			return
		}

		s.recordersMutex.RLock()
		defer s.recordersMutex.RUnlock()

		for recorder := range s.recorders[sensorType] {
			s.wg.Add(1)
			go func(rec Recorder, meas *sensor.Measurement) {
				if err := rec.Record(meas); err != nil {
					logutil.Warn("Error recording measurement of type %s with recorder %v: %v", sensorType, rec, err)
				}
				s.wg.Done()
			}(recorder, measurement)
		}
	})
}

// Close closes all recorders
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

// Start starts subscriptions for all recorders
func (s *Manager) Start() error {
	s.recordersMutex.RLock()
	defer s.recordersMutex.RUnlock()

	for sensorType, recorders := range s.recorders {
		for _, qos := range recorders {
			if err := s.SubscribeToSensor(sensorType, qos); err != nil {
				return fmt.Errorf("Failed to subscribe to sensor type %s: %v", sensorType, err)
			}
		}
	}
	return nil
}
