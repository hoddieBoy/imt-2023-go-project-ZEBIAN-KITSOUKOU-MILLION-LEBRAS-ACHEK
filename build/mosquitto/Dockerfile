FROM eclipse-mosquitto:2.0.18

ARG MQTT_USER
ARG MQTT_PASSWORD

# Set the configuration file
COPY conf/mosquitto.conf /mosquitto/config/mosquitto.conf

# Set the password file
RUN echo "$MQTT_USER:$MQTT_PASSWORD" > /mosquitto/config/passwd

# Encrypt the password file
RUN chmod 600 /mosquitto/config/passwd
RUN mosquitto_passwd -U /mosquitto/config/passwd
