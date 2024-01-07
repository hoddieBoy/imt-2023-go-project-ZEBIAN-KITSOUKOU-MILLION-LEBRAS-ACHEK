package mqtt_helper

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	client mqtt.Client
	config *MQTTConfig
}

func NewClient(config *MQTTConfig, clientID string) *MQTTClient {
	opts := mqtt.NewClientOptions().AddBroker(config.GetServerAddress()).SetClientID(clientID)

	if config.Server.Username != "" {
		opts.SetUsername(config.Server.Username)
	}

	if config.Server.Password != "" {
		opts.SetPassword(config.Server.Password)
	}

	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		fmt.Printf("\033[33mconnection lost, trying to reconnect:\n\t<<%s>>\033[0m\n", err)
		if token := client.Connect(); token.Wait() && token.Error() != nil {
			fmt.Printf("\033[31mfailed to reconnect:\n\t<<%s>>\033[0m\n", token.Error())
		}
	})

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		fmt.Printf("\033[32mconnection established with broker %s\033[0m\n", config.GetServerAddress())
	})

	return &MQTTClient{
		client: mqtt.NewClient(opts),
		config: config,
	}
}

func (c *MQTTClient) Connect() error {
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("\033[31mfailed to connect to broker:\n\t<<%w>>\033[0m", token.Error())
	}
	return nil
}

func (c *MQTTClient) Disconnect() {
	c.client.Disconnect(250)
	fmt.Printf("\033[32mdisconnected from broker %s\033[0m\n", c.config.GetServerAddress())
}

func (c *MQTTClient) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error {
	if token := c.client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		return fmt.Errorf("\033[31mfailed to subscribe to topic %s:\n\t<<%w>>\033[0m", topic, token.Error())
	}

	return nil
}

func (c *MQTTClient) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	token := c.client.Publish(topic, qos, retained, payload)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("\033[31mfailed to publish to topic %s:\n\t<<%w>>\033[0m", topic, token.Error())
	}
	return nil
}

func (c *MQTTClient) Unsubscribe(topics ...string) error {
	if token := c.client.Unsubscribe(topics...); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to unsubscribe from topics %v:\n\t<<%w>>", topics, token.Error())
	}
	return nil
}
