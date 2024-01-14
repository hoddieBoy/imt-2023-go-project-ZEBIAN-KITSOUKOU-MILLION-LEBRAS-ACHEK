package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"imt-atlantique.project.group.fr/meteo-airport/internal/log"
)

type Client struct {
	client mqtt.Client
	config *Config
}

func NewClient(config *Config) *Client {
	opts := mqtt.NewClientOptions().AddBroker(config.GetServerAddress()).SetClientID(config.ClientID)

	if config.Username != "" {
		opts.SetUsername(config.Username)
	}

	if config.Password != "" {
		opts.SetPassword(config.Password)
	}

	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Warn("connection lost with broker %s", config.GetServerAddress())
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			log.Error("failed to reconnect to broker %s:\n\t<<%v>>", config.GetServerAddress(), token.Error())
		}
	})

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Info("connected to broker %s", config.GetServerAddress())
	})

	return &Client{
		client: mqtt.NewClient(opts),
		config: config,
	}
}

func (c *Client) Connect() error {
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		log.Error("failed to connect to broker %s:\n\t<<%v>>", c.config.GetServerAddress(), token.Error())
		return token.Error()
	}

	return nil
}

func (c *Client) Disconnect() {
	c.client.Disconnect(250)
	log.Info("disconnected from broker %s", c.config.GetServerAddress())
}

func (c *Client) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error {
	if token := c.client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		log.Error("failed to subscribe to topic %s:\n\t<<%v>>", topic, token.Error())
		return token.Error()
	}

	return nil
}

func (c *Client) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	token := c.client.Publish(topic, qos, retained, payload)
	if token.Wait() && token.Error() != nil {
		log.Error("failed to publish to topic %s:\n\t<<%v>>", topic, token.Error())
		return token.Error()
	}

	return nil
}

func (c *Client) Unsubscribe(topics ...string) error {
	if token := c.client.Unsubscribe(topics...); token.Wait() && token.Error() != nil {
		log.Error("failed to unsubscribe from topics %v:\n\t<<%v>>", topics, token.Error())
		return token.Error()
	}

	return nil
}
