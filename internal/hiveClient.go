package internal

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
)

/*
	client := makeDefaultClient()

	subscribeWithQos_1(client, "topic/test")
	publish(client, "This is an example", "topic/test")

	time.Sleep(time.Second)

	client.Disconnect(250)
*/

func CreateFileLogger(filePath string) (*log.Logger, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	logger := log.New(file, "FILE: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.SetOutput(file)
	return logger, nil
}

var log1, _ = CreateFileLogger("logs.txt")

func makeDefaultClient() mqtt.Client {
	var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		log1.Printf("Received message \"%s\" from topic \"%s\"\n", msg.Payload(), msg.Topic())
	}

	var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
		log1.Println("Connected")
	}

	var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
		log1.Printf("Connection lost: %v", err)
	}
	var broker = "411d6c045163486b846c891f3910e83f.s2.eu.hivemq.cloud"
	var port = 8883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tls://%s:%d", broker, port))
	opts.SetClientID("aClientId")
	opts.SetUsername("arduino")
	opts.SetPassword("not4nArduino")

	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}

func subscribeWithQos_1(client mqtt.Client, topic string) {
	token := client.Subscribe(topic, 1, nil)

	token.Wait()
	if token.Error() != nil {
		log1.Printf("Failed to subscribe to topic")
		panic(token.Error())
	}
	log1.Printf("Subscribed to topic: %s", topic)
}

func publish(client mqtt.Client, message string, topic string) {
	token := client.Publish(topic, 0, false, message)
	token.Wait()
	if token.Error() != nil {
		log1.Printf("Failed to publish to topic")
		panic(token.Error())
	}
	log1.Printf("Sent \"%s\" to topic \"%s\"", message, topic)
}
