package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	ABB "github.com/rpoisel/modbus-mqtt/abb"
)

func setupMqtt() *MQTT.ClientOptions {
	opts := MQTT.NewClientOptions()
	opts.AddBroker("ssl://hostname.tld:8883")
	opts.SetClientID("ssl-sample")
	opts.SetUsername("user")
	opts.SetPassword("pass")
	return opts
}

const (
	solarPowerID    = 1
	obtainedPowerID = 2
)

func main() {
	powerMeters := make(map[byte]*ABB.B23)
	for _, id := range []byte{obtainedPowerID, solarPowerID} {
		b23Instance, err := ABB.NewB23("/dev/ttyUSB0", id)
		if err != nil {
			panic(err.Error())
		}
		defer b23Instance.Close()
		powerMeters[id] = b23Instance
	}

	mqttClient := MQTT.NewClient(setupMqtt())
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer mqttClient.Disconnect(250)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	for {
		select {
		case <-stopChan:
			return
		default:
			obtainedPower, err := powerMeters[obtainedPowerID].Power()
			if err != nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			time.Sleep(100 * time.Millisecond)

			solarPower, err := powerMeters[solarPowerID].Power()
			if err != nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			time.Sleep(1000 * time.Millisecond)

			var totalPower int32
			if solarPower > 0 {
				totalPower = solarPower + obtainedPower
			} else {
				totalPower = obtainedPower
			}

			text := fmt.Sprintf("%d", obtainedPower)
			mqttClient.Publish("/homeautomation/power/obtained", 0, false, text)
			text = fmt.Sprintf("%d", solarPower)
			mqttClient.Publish("/homeautomation/power/solar", 0, false, text)
			text = fmt.Sprintf("%d", totalPower)
			mqttClient.Publish("/homeautomation/power/total", 0, false, text)

			buf := new(bytes.Buffer)
			binary.Write(buf, binary.LittleEndian, solarPower)
			binary.Write(buf, binary.LittleEndian, obtainedPower)
			binary.Write(buf, binary.LittleEndian, totalPower)
			mqttClient.Publish("/homeautomation/power/cumulative", 0, false, buf.Bytes())
		}
	}
}
