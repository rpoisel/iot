package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	UTIL "github.com/rpoisel/IoT/internal/util"
)

const (
	solarPowerID    = 1
	obtainedPowerID = 2
)

type configuration struct {
	Mqtt   UTIL.MqttConfiguration
	Modbus struct {
		Device string
	}
}

func main() {
	configuration := configuration{}
	UTIL.ReadConfig("/etc/homeautomation.json", &configuration)

	powerMeters := make(map[byte]*B23)
	for _, id := range []byte{obtainedPowerID, solarPowerID} {
		b23Instance, err := NewB23(configuration.Modbus.Device, id)
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
			var err error
			var readings UTIL.Readings

			readings.Obtained, err = powerMeters[obtainedPowerID].Power()
			if err != nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			time.Sleep(100 * time.Millisecond)

			readings.Solar, err = powerMeters[solarPowerID].Power()
			if err != nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			time.Sleep(1000 * time.Millisecond)

			if readings.Solar > 0 {
				readings.Total = readings.Solar + readings.Obtained
			} else {
				readings.Total = readings.Obtained
			}

			text := fmt.Sprintf("%d", readings.Obtained)
			mqttClient.Publish("/homeautomation/power/obtained", 0, false, text)
			text = fmt.Sprintf("%d", readings.Solar)
			mqttClient.Publish("/homeautomation/power/solar", 0, false, text)
			text = fmt.Sprintf("%d", readings.Total)
			mqttClient.Publish("/homeautomation/power/total", 0, false, text)

			mqttClient.Publish("/homeautomation/power/cumulative", 0, false, readings.ToBuf())
		}
	}
}
