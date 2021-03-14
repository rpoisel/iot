package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	MODBUS "github.com/goburrow/modbus"
	UTIL "github.com/rpoisel/iot/internal/util"
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
	var configPath = flag.String("c", "/etc/homeautomation.yaml", "Path to the configuration file")
	flag.Parse()

	configuration := configuration{}
	UTIL.ReadConfig(*configPath, &configuration)

	powerMeters := make(map[byte]*b23)
	for _, id := range []byte{obtainedPowerID, solarPowerID} {
		b23Instance, err := newB23(configuration.Modbus.Device, id)
		if err != nil {
			panic(err.Error())
		}
		defer b23Instance.Close()
		powerMeters[id] = b23Instance
	}

	mqttClient := UTIL.SetupMqtt(configuration.Mqtt, nil, nil)
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

type b23 struct {
	handler *MODBUS.RTUClientHandler
	client  MODBUS.Client
}

func newB23(device string, slaveID byte) (b23Instance *b23, err error) {
	b23Instance = &b23{}
	b23Instance.handler = MODBUS.NewRTUClientHandler(device)
	b23Instance.handler.BaudRate = 19200
	b23Instance.handler.DataBits = 8
	b23Instance.handler.Parity = "E"
	b23Instance.handler.StopBits = 1
	b23Instance.handler.SlaveId = slaveID
	// b23Instance.handler.Logger = log.New(os.Stdout, "rtu: ", log.LstdFlags)
	b23Instance.handler.Logger = nil
	err = b23Instance.handler.Connect()
	if err != nil {
		return
	}
	b23Instance.client = MODBUS.NewClient(b23Instance.handler)
	return
}

func (b23Instance *b23) Power() (power int32, err error) {
	var rawData []byte
	rawData, err = b23Instance.client.ReadHoldingRegisters(0x5B00, 66)
	if err != nil || rawData == nil {
		return
	}
	buf := bytes.NewReader(rawData[40:44])
	err = binary.Read(buf, binary.BigEndian, &power)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	power /= 100
	return
}

func (b23Instance *b23) Close() {
	b23Instance.handler.Close()
}
