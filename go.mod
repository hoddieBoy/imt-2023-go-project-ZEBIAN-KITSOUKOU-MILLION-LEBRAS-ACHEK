module imt-atlantique.project.group.fr/meteo-airport

go 1.21.4

require (
	github.com/eclipse/paho.mqtt.golang v1.4.3
	github.com/fatih/color v1.16.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
)

replace internal => ./internal
