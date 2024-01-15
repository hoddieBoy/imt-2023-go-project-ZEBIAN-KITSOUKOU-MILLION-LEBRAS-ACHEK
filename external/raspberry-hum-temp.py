#!/usr/bin/env python3

import RPi.GPIO as GPIO
import time
import Freenove_DHT as DHT
import paho.mqtt.client as mqtt
from datetime import datetime
from enum import Enum
import json

class MeasurementType(Enum):
    TEMPERATURE = 1
    HUMIDITY = 2

class Measurement:
    def __init__(self, sensor_id: int, airport_id: str, type: MeasurementType, value: float, unit: str, timestamp: datetime):
        self.sensor_id = sensor_id
        self.airport_id = airport_id
        self.type = type
        self.value = value
        self.unit = unit
        self.timestamp = timestamp

class MeasurementEncoder(json.JSONEncoder):
    def default(self, obj):
        if isinstance(obj, Measurement):
            return {
                'sensor_id': obj.sensor_id,
                'airport_id': obj.airport_id,
                'type': obj.type.name,  # Convert Enum to string
                'value': obj.value,
                'unit': obj.unit,
                'timestamp': obj.timestamp.isoformat()  # Convert datetime to ISO 8601 string
            }
        return super().default(obj)
    


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
        
        measurement = Measurement(1, 'airport1', MeasurementType.TEMPERATURE, dht.temperature, 'C', datetime.now())
        json_str = json.dumps(measurement, cls=MeasurementEncoder)
        
        
        client.publish(topic, f"{json_str}")
        time.sleep(5)       
        
if __name__ == '__main__':
    print ('Program is starting ... ')
    try:
        loop()
    except KeyboardInterrupt:
        GPIO.cleanup()
        exit()  

