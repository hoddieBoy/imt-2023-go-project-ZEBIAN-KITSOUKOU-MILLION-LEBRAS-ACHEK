# Projet Weather Data Collection and Retrieval

## Description

Meteo Airport is a comprehensive project that simulates a weather station at an airport. It collects data from various
sensors, including temperature, humidity, wind speed, and pressure. The data is then published to an MQTT broker and
stored in InfluxDB. The project also includes an alert manager that listens to the MQTT broker and triggers alerts based
on predefined conditions.

## Features

- Sensor data simulation for temperature, humidity, wind speed, and pressure
- MQTT publisher for sensor data
- Storage manager for storing sensor data in InfluxDB or in CSV files
- Alert manager for triggering alerts based on predefined conditions
- REST API for retrieving sensor data and alerts
- Swagger documentation for the REST API
- Test suite for validating the project
- Linter for ensuring code quality
- Docker Compose for running services locally
- Makefile for automating common tasks
- CI/CD pipeline for automating the build and test process

## Project structure

The project is structured into several packages, each responsible for a specific functionality:

- `cmd/` contains the mains applications(Storage manager, Alert manager, Sensor data simulation, REST API)
- `internal/` contains the internal packages used by the applications
- `config/` contains the configuration files for the applications
- `scripts/` contains the scripts for building, testing, and other operations
- `test/` contains the tests for the project
- `external/` contains the external packages used by the project

## Getting started

### Prerequisites

- Go 1.21.4 or compatible
- Make 3.81 or compatible
- Docker 24.0.6 or compatible
- An MQTT broker(
  e.g. [Mosquitto](https://mosquitto.org/), [HiveMQ](https://www.hivemq.com/), [EMQ X](https://www.emqx.io/))
- An InfluxDB database)

### Installation steps

1. Clone the repository:
    ```bash
    git clone https://github.com/jarhead-killgrave/imt-2023-go-project-ZEBIAN-KITSOUKOU-MILLION-LEBRAS-ACHEK.git
    ```

2. Navigate to the project directory:
    ```bash
    cd imt-2023-go-project-ZEBIAN-KITSOUKOU-MILLION-LEBRAS-ACHEK
    ```

3. Copy the `.env.dev` file to `.env`, and edit it to match your configuration:
   ```bash
   cp .env.dev .env
   ```
   or for Windows:
   ```bash
    copy .env.dev .env
   ```

4. Build the applications:
   ```bash
   make build
   ```

5. Initialize all necessary services:
   ```bash
   make init
   ```
   Alternatively, you can initialize each service without using `make`:
   ```bash
   docker compose up -d
   ```
   After initialization, you should be able to access the following services:
	- [http://localhost:8086](http://localhost:8086) for InfluxDB
	- [http://localhost:1883](http://localhost:1883) for the MQTT broker

6. For running a specific application, you can use the following commands:
   ```bash
   ./<application_name> <configuration_file>
   ```
   For example, to run the storage manager:
   ```bash
   ./storage-manager config/storage-manager.yaml
   ```
   or for Windows:
   ```bash
   .\<application_name>.exe <configuration_file>
   ```
   For example, to run the storage manager:
   ```bash
   .\storage-manager.exe config\storage-manager.yaml
   ```

You can run all the applications at once by running the following command:

```bash
make run
```

# Configuration

The project uses YAML files for configuration. You can find the configuration files in the `config/` directory.
Each application has its own configuration file. The configuration contains all the necessary information for running
the application, including the MQTT broker address, the InfluxDB address, the topic names...

# Testing

The project includes a test suite for validating the project. The tests are located in the `test/` directory.
You can run the tests with the following command:

```bash
make test
```

or with command line:

```bash
go test ./test/...
```

# Linter

Before your branch is merged, [golangci-lint](https://golangci-lint.run/) will be run on your code on the CI server.

First, you need to install it locally:

```bash
docker pull golangci/golangci-lint
```

After You can run it locally with docker by running the following command:

```bash
docker run -t --rm -v $(pwd):/app -w /app golangci/golangci-lint golangci-lint run -v
```

or with make:

```bash
make lint
```

If you want to fix the issues automatically, you can run the following command:

```bash
docker run -t --rm -v $(pwd):/app -w /app golangci/golangci-lint golangci-lint run -v --fix
```

or with make:

```bash
make lint-fix
```

## Authors

- [ZEBIAN Jana](https://github.com/JanaZebian)
- [MILLION Julien](https://github.com/AlphaOrOmega)
- [LEBRAS Gregoire](https://github.com/gregoireLeBras)
- [KITSOUKOU Manne Émile](https://github.com/jarhead-killgrave)
- [ACHEK Jamil](https://github.com/JamWare)
