package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	UTIL "github.com/rpoisel/IoT/internal/util"
)

type loxoneConfiguration struct {
	Miniserver string
	Username   string
	Password   string
	MqttPath   string
	Blinds     map[string]string
}

type configuration struct {
	Mqtt   UTIL.MqttConfiguration
	Loxone loxoneConfiguration
}

func sendHTTPGetRequest(path string, username string, password string) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", path, nil)
	req.SetBasicAuth(username, password)
	client.Do(req)
}

func blindsPublishHandler(_ MQTT.Client, msg MQTT.Message) {
	srcBlind := strings.Replace(string(msg.Topic()), "/homeautomation/blinds/", "", -1)
	loxoneBlind, exists := config.Loxone.Blinds[srcBlind]
	if !exists {
		log.Println("Blind does not exist: ", srcBlind)
		return
	}

	url := "http://" + config.Loxone.Miniserver + "/dev/sps/io/" + loxoneBlind + "/"
	payload := strings.ToLower(string(msg.Payload()))
	if payload == "up" {
		url += "Up"
	} else if payload == "down" {
		url += "Down"
	} else {
		return
	}
	go sendHTTPGetRequest(url, config.Loxone.Username, config.Loxone.Password)
}

var config configuration = configuration{}

func main() {
	var configPath = flag.String("c", "/etc/homeautomation.yaml", "Path to the configuration file")
	flag.Parse()

	UTIL.ReadConfig(*configPath, &config)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	mqttClient := UTIL.SetupMqtt(config.Mqtt, func(_ MQTT.Client, msg MQTT.Message) {
		log.Print("Unhandled MQTT message ", msg)
	}, func(client MQTT.Client) {
		for src := range config.Loxone.Blinds {
			client.Subscribe("/homeautomation/blinds/"+src, 0 /* qos */, blindsPublishHandler)
		}
	})
	defer mqttClient.Disconnect(250)

	<-stopChan
	fmt.Println("Good bye!")
}
