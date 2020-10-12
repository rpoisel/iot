package main

import (
	"bytes"
	"encoding/binary"
	"fmt"

	MODBUS "github.com/goburrow/modbus"
)

type B23 struct {
	handler *MODBUS.RTUClientHandler
	client  MODBUS.Client
}

func NewB23(device string, slaveID byte) (b23 *B23, err error) {
	b23 = &B23{}
	b23.handler = MODBUS.NewRTUClientHandler(device)
	b23.handler.BaudRate = 19200
	b23.handler.DataBits = 8
	b23.handler.Parity = "E"
	b23.handler.StopBits = 1
	b23.handler.SlaveId = slaveID
	// b23.handler.Logger = log.New(os.Stdout, "rtu: ", log.LstdFlags)
	b23.handler.Logger = nil
	err = b23.handler.Connect()
	if err != nil {
		return
	}
	b23.client = MODBUS.NewClient(b23.handler)
	return
}

func (b23 *B23) Power() (power int32, err error) {
	var rawData []byte
	rawData, err = b23.client.ReadHoldingRegisters(0x5B00, 66)
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

func (b23 *B23) Close() {
	b23.handler.Close()
}
