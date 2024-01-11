package mqtt_helper

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"imt-atlantique.project.group.fr/meteo-airport/internal/logutil"
)

type MQTTClient struct {
	client mqtt.Client
	config *MQTTConfig
}

func NewClient(config *MQTTConfig, clientID string) *MQTTClient {
	opts := mqtt.NewClientOptions().AddBroker(config.GetServerAddress()).SetClientID(clientID)

	if config.Username != "" {
		opts.SetUsername(config.Username)
	}

	if config.Password != "" {
		opts.SetPassword(config.Password)
	}

	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		logutil.Warn("connection lost with broker %s", config.GetServerAddress())
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			logutil.Error("failed to reconnect to broker %s:\n\t<<%v>>", config.GetServerAddress(), token.Error())
		}
	})

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		logutil.Info("connected to broker %s", config.GetServerAddress())
	})

	return &MQTTClient{
		client: mqtt.NewClient(opts),
		config: config,
	}
}

func (c *MQTTClient) Connect() error {
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		logutil.Error("failed to connect to broker %s:\n\t<<%v>>", c.config.GetServerAddress(), token.Error())
		return token.Error()
	}
	return nil
}

func (c *MQTTClient) Disconnect() {
	c.client.Disconnect(250)
	logutil.Info("disconnected from broker %s", c.config.GetServerAddress())
}

func (c *MQTTClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error {
	if token := c.client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		logutil.Error("failed to subscribe to topic %s:\n\t<<%v>>", topic, token.Error())
		return token.Error()
	}

	return nil
}

func (c *MQTTClient) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	token := c.client.Publish(topic, qos, retained, payload)
	if token.Wait() && token.Error() != nil {
		logutil.Error("failed to publish to topic %s:\n\t<<%v>>", topic, token.Error())
		return token.Error()
	}
	return nil
}

func (c *MQTTClient) Unsubscribe(topics ...string) error {
	if token := c.client.Unsubscribe(topics...); token.Wait() && token.Error() != nil {
		logutil.Error("failed to unsubscribe from topics %v:\n\t<<%v>>", topics, token.Error())
		return token.Error()
	}
	return nil
}
