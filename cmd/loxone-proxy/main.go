package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	UTIL "github.com/rpoisel/IoT/internal/util"
)

type LoxoneConfiguration struct {
	Miniserver string
	Username   string
	Password   string
	MqttPath   string
	Blinds     map[string]string
}

type Configuration struct {
	Mqtt   UTIL.MqttConfiguration
	Loxone LoxoneConfiguration
}

func sendHttpGetRequest(path string, username string, password string) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", path, nil)
	req.SetBasicAuth(username, password)
	client.Do(req)
}

func defaultMqttPublishHandler(_ MQTT.Client, msg MQTT.Message) {
	log.Print("Unhandled MQTT message ", msg)
}

func blindsPublishHandler(_ MQTT.Client, msg MQTT.Message) {
	srcBlind := strings.Replace(string(msg.Topic()), "/homeautomation/blinds/", "", -1)
	loxoneBlind, exists := configuration.Loxone.Blinds[srcBlind]
	if !exists {
		log.Println("Blind does not exist: ", srcBlind)
		return
	}

	url := "http://" + configuration.Loxone.Miniserver + "/dev/sps/io/" + loxoneBlind + "/"
	payload := strings.ToLower(string(msg.Payload()))
	if payload == "up" {
		url += "Up"
	} else if payload == "down" {
		url += "Down"
	} else {
		return
	}
	go sendHttpGetRequest(url, configuration.Loxone.Username, configuration.Loxone.Password)
}

var configuration Configuration = Configuration{}

func main() {
	UTIL.ReadConfig("/etc/homeautomation.json", &configuration)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	mqttClient := UTIL.SetupMqtt(configuration.Mqtt, defaultMqttPublishHandler)
	defer mqttClient.Disconnect(250)

	for src := range configuration.Loxone.Blinds {
		mqttClient.Subscribe("/homeautomation/blinds/"+src, 0 /* qos */, blindsPublishHandler)
	}

	<-stopChan
	fmt.Println("Good bye!")
}
