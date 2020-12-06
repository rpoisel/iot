package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	UTIL "github.com/rpoisel/iot/internal/util"
	JoySticks "github.com/splace/joysticks"
)

type configuration struct {
	Mqtt UTIL.MqttConfiguration
}

func main() {
	var configPath = flag.String("c", "/etc/homeautomation.yaml", "Path to the configuration file")
	flag.Parse()

	configuration := configuration{}
	UTIL.ReadConfig(*configPath, &configuration)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	mqttClient := UTIL.SetupMqtt(configuration.Mqtt, nil, nil)
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
