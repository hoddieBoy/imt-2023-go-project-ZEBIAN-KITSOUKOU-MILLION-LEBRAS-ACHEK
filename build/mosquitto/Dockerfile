FROM eclipse-mosquitto:2.0.18

ARG MQTT_USER
ARG MQTT_PASSWORD

# Set the configuration file
COPY conf/mosquitto.conf /mosquitto/config/mosquitto.conf

# Create the password file
RUN touch /mosquitto/config/passwd
RUN chown 600 /mosquitto/config/passwd
RUN mosquitto_passwd -b /mosquitto/config/passwd "$MQTT_USER" "$MQTT_PASSWORD"
