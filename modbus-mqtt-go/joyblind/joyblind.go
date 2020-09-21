package main

import (
	"log"
	"os"
	"os/signal"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	CONF "github.com/rpoisel/modbus-mqtt/conf"
	JoySticks "github.com/splace/joysticks"
)

func setupMqtt(conf *CONF.Config) *MQTT.ClientOptions {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(conf.ValueAsString("broker"))
	opts.SetUsername(conf.ValueAsString("username"))
	opts.SetPassword(conf.ValueAsString("password"))
	return opts
}

func main() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	conf, err := CONF.NewConfig("/etc/homeautomation.json")
	if err != nil {
		panic(err)
	}

	mqttClient := MQTT.NewClient(setupMqtt(conf.Value("mqtt")))
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
