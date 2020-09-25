package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	JoySticks "github.com/splace/joysticks"
)

type MqttConfiguration struct {
	Broker   string
	Username string
	Password string
}

type Configuration struct {
	Mqtt MqttConfiguration
}

func setupMqtt(config MqttConfiguration) (opts *MQTT.ClientOptions) {
	opts = MQTT.NewClientOptions()
	opts.AddBroker(config.Broker)
	opts.SetUsername(config.Username)
	opts.SetPassword(config.Password)
	return opts
}

func readConfig(path string) (config *Configuration) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		panic(err)
	}
	return &configuration
}

func main() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	configuration := readConfig("/etc/homeautomation.json")

	mqttClient := MQTT.NewClient(setupMqtt(configuration.Mqtt))
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer mqttClient.Disconnect(250)

	evts := JoySticks.Capture(
		JoySticks.Channel{1, JoySticks.HID.OnClose},
		JoySticks.Channel{2, JoySticks.HID.OnClose},
		JoySticks.Channel{1, JoySticks.HID.OnMove},
	)

	for {
		select {
		case <-stopChan:
			log.Print("Gracefully shutting down ...")
			return
		case <-evts[0]:
			mqttClient.Publish("/homeautomation/blinds/SR", 2, false, "up")
		case <-evts[1]:
			mqttClient.Publish("/homeautomation/blinds/SR", 2, false, "down")
		case h := <-evts[2]:
			hpos, ok := h.(JoySticks.CoordsEvent)
			if !ok {
				return
			}
			if hpos.Y == -1 {
				mqttClient.Publish("/homeautomation/blinds/SR", 2, false, "up")
			} else if hpos.Y == 1 {
				mqttClient.Publish("/homeautomation/blinds/SR", 2, false, "down")
			}
		}
	}
}
