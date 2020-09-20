package main

import (
	"log"
	"os"
	"os/signal"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	CONF "github.com/rpoisel/modbus-mqtt/conf"
	JoySticks "github.com/splace/joysticks"
)

func setupMqtt() *MQTT.ClientOptions {
	config, err := CONF.ReadConfigSection("/etc/homeautomation.json", "mqtt")
	if err != nil {
		panic(err)
	}

	opts := MQTT.NewClientOptions()
	opts.AddBroker(config["broker"].(string))
	opts.SetUsername(config["username"].(string))
	opts.SetPassword(config["password"].(string))
	return opts
}

func main() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	mqttClient := MQTT.NewClient(setupMqtt())
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
