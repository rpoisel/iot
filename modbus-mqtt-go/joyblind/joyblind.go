package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	JoySticks "github.com/splace/joysticks"
)

func setupMqtt() *MQTT.ClientOptions {
	config, err := readConfigSection("/etc/homeautomation.json", "mqtt")
	if err != nil {
		panic(err)
	}

	opts := MQTT.NewClientOptions()
	opts.AddBroker(config["broker"].(string))
	opts.SetUsername(config["username"].(string))
	opts.SetPassword(config["password"].(string))
	return opts
}

func readConfigSection(path string, section string) (resultMap map[string]interface{}, err error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var dat map[string]interface{}
	if err := json.Unmarshal(byteValue, &dat); err != nil {
		return nil, err
	}
	resultMap = dat[section].(map[string]interface{})
	return
}

func main() {
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
