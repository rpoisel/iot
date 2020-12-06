package main

import (
	"bytes"
	"encoding/binary"
	"fmt"

	MODBUS "github.com/goburrow/modbus"
)

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
