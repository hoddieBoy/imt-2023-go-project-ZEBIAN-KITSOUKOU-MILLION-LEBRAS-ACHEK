#!/usr/bin/env python3

import RPi.GPIO as GPIO
import time
import Freenove_DHT as DHT
import paho.mqtt.client as mqtt
from datetime import datetime
from enum import Enum
import json

class MeasurementType(Enum):
    Temperature = 1
    Humidity = 2

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
topic = "airport/T"
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

adc = ADCDevice() # Define an ADCDevice class object

def setup():
    global adc
    if(adc.detectI2C(0x48)): #Detect the pcf8591.
        adc = PCF8591()
    elif(adc.detectI2C(0x4b)): #Detect the ads7830.
        adc = ADS7830()
    else:
        print("No correct I2C address found, \n"
        "Please use command 'i2cdetect -y 1' to check the I2C address! \n"
        "Program Exit. \n")
        exit(-1)

def loop():
    while(True):
        value = adc.analogRead(0) #read the ADC value of A0 pin     
        voltage = value / 255.0 * 3.3 #calculate the voltage value
        Rt = 10 * voltage / (3.3 - voltage) #calculate the resistance value of the thermistor
        tempK = 1/(1/(273.15 + 25) + math.log(Rt/10)/3950.0) #calculate the temperature(Kelvin)
        tempC = tempK -273.15 #calculate the temperature(Celsius)
        print ('ADC Value : %d, Voltage : %.2f, Temperature : %.2f'%(value,voltage,tempC))
        measurement = Measurement(1, 'AKV', MeasurementType.Temperature, str(tempC), 'Celsius', datetime.now())
        json_str = json.dumps(measurement, cls=MeasurementEncoder)
        
        client.publish(topic, f"{json_str}")
        time.sleep(3)
        
if __name__ == '__main__':
    print ('Program is starting ... ')
    try:
        loop()
    except KeyboardInterrupt:
        GPIO.cleanup()
        exit()  

