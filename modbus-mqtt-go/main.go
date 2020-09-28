package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"os/signal"
	"time"

	ABB "github.com/rpoisel/modbus-mqtt/abb"
	UTIL "github.com/rpoisel/modbus-mqtt/util"
)

const (
	solarPowerID    = 1
	obtainedPowerID = 2
)

type Configuration struct {
	Mqtt   UTIL.MqttConfiguration
	Modbus struct {
		Device string
	}
}

func main() {
	configuration := Configuration{}
	UTIL.ReadConfig("/etc/homeautomation.json", &configuration)

	powerMeters := make(map[byte]*ABB.B23)
	for _, id := range []byte{obtainedPowerID, solarPowerID} {
		b23Instance, err := ABB.NewB23(configuration.Modbus.Device, id)
		if err != nil {
			panic(err.Error())
		}
		defer b23Instance.Close()
		powerMeters[id] = b23Instance
	}

	mqttClient := UTIL.SetupMqtt(configuration.Mqtt, nil)
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
