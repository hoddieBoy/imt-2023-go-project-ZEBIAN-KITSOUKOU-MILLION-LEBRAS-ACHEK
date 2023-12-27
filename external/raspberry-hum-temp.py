#!/usr/bin/env python3

import RPi.GPIO as GPIO
import time
import Freenove_DHT as DHT
import paho.mqtt.client as mqtt

# Define MQTT broker and topic
broker = "localhost"
topic = "capteur/TH"
port = 1883

id = "inform your ID"
passwd = "inform your password"

# Define callback functions
def on_connect(client, userdata, flags, rc):
    print("Connected to MQTT broker")
    client.subscribe(topic)

def on_message(client, userdata, msg):
    print(f"Received message: {msg.payload.decode()}")

# Create MQTT client instance
client = mqtt.Client(id, clean_session=False)

client.username_pw_set(id, passwd)

# Set callback functions
client.on_connect = on_connect
client.on_message = on_message

# Connect to MQTT broker
client.connect(broker, port)

# Start the MQTT loop
client.loop_start()

DHTPin = 11     #define the pin of DHT11

def loop():
    dht = DHT.DHT(DHTPin)   #create a DHT class object
    counts = 0 # Measurement counts
    while(True):
        counts += 1
        print("Measurement counts: ", counts)
        for i in range(0,15):            
            chk = dht.readDHT11()     #read DHT11 and get a return value. Then determine whether data read is normal according to the return value.
            if (chk is dht.DHTLIB_OK):      #read DHT11 and get a return value. Then determine whether data read is normal according to the return value.
                print("DHT11,OK!")
                break
            time.sleep(0.1)
        print("Humidity : %.2f, \t Temperature : %.2f \n"%(dht.humidity,dht.temperature))
        client.publish(topic, f"Temperature: {dht.temperature}, Humidity: {dht.humidity}")
        time.sleep(5)       
        
if __name__ == '__main__':
    print ('Program is starting ... ')
    try:
        loop()
    except KeyboardInterrupt:
        GPIO.cleanup()
        exit()  

